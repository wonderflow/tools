package pick

import (
	"fmt"
	"sort"
	"time"
)

type InfluxD struct {
	Id               string    `json:"id" bson:"id"` //host:port_diskTag
	Dirs             []string  `json:"dirs" bson:"dirs"`
	Host             string    `json:"host" bson:"host"`
	Port             string    `json:"port" bson:"port"`
	SpaceUsed        float64   `json:"spaceUsed" bson:"spaceused"`
	SpaceFree        float64   `json:"spaceFree" bson:"spacefree"`
	SpaceConsumeRate float64   `json:"spaceConsumeRate" bson:"spaceconsumerate"`
	Group            string    `json:"group" bson:"group"`
	LastTouch        time.Time `json:"lastTouch" bson:"lasttouch"`
}

type InfluxDs []*InfluxD

func (a InfluxDs) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a InfluxDs) Len() int      { return len(a) }
func (a InfluxDs) Less(i, j int) bool {
	if a[i].Group == "" && a[j].Group == "" {
		if a[i].SpaceFree == a[j].SpaceFree {
			if a[i].Host == a[j].Host {
				return a[i].Port < a[j].Port
			} else {
				return a[i].Host < a[j].Host
			}
		} else {
			return a[i].SpaceFree > a[j].SpaceFree
		}
	} else if a[i].Group == "" {
		return true
	} else if a[j].Group == "" {
		return false
	}

	if a[i].SpaceFree == a[j].SpaceFree {
		if a[i].Host == a[j].Host {
			return a[i].Port < a[j].Port
		} else {
			return a[i].Host < a[j].Host
		}
	} else {
		return a[i].SpaceFree > a[j].SpaceFree
	}
}

func search(influxds InfluxDs, searchindex int, cur_result *[]int, cur_index int, cursum float64, final_result *[]int, finalsum *float64, replica int, machinemap map[string]bool, tolerate float64) {
	if cur_index == replica {
		if *finalsum == -1 || cursum < *finalsum {
			*finalsum = cursum
			*final_result = *cur_result
		}
		return
	}
	for i := searchindex; i < len(influxds); i++ {
		if influxds[i].Group != "" {
			continue
		}
		if value := machinemap[influxds[i].Host]; !value {
			tmpsum := 0.0
			if cur_index > 0 {
				tmpsum += influxds[(*cur_result)[cur_index-1]].SpaceFree - influxds[i].SpaceFree
			}
			if *finalsum < 0 || cursum+tmpsum < *finalsum {
				machinemap[influxds[i].Host] = true
				(*cur_result)[cur_index] = i
				search(influxds, i+1, cur_result, cur_index+1, cursum+tmpsum, final_result, finalsum, replica, machinemap, tolerate)
				if *finalsum >= 0 && *finalsum <= tolerate {
					return
				}
				machinemap[influxds[i].Host] = false
			}
		}
	}
}

func PickInfluxds(influxds InfluxDs, replica int) (ns []*InfluxD, err error) {
	if len(influxds) < replica {
		err = fmt.Errorf("No enough influxdb to pick.")
		return
	}
	machinemap := make(map[string]bool)
	cur_result := make([]int, replica)
	final_result := make([]int, replica)
	sort.Sort(influxds)

	finalsum := -1.0
	search(influxds, 0, &cur_result, 0, 0, &final_result, &finalsum, replica, machinemap, 10.0)
	if finalsum < 0 {
		ns = influxds[0:replica]
		return
	}
	for i := range final_result {
		ns = append(ns, influxds[final_result[i]])
	}
	return
}
