package pick

import (
	"github.com/davecgh/go-spew/spew"
	"reflect"
	"testing"
)

func color(ns []*InfluxD) {
	for i := range ns {
		ns[i].Group = "colored"
	}
}

func Test_MakeGroup(t *testing.T) {

	influxs := [][]*InfluxD{
		{
			&InfluxD{Host: "1one.pandora.com", Port: "08086", SpaceUsed: 0, SpaceFree: 1024 * 10, SpaceConsumeRate: 1024},
			&InfluxD{Host: "1one.pandora.com", Port: "09086", SpaceUsed: 0, SpaceFree: 1024 * 8, SpaceConsumeRate: 1024},
			&InfluxD{Host: "1one.pandora.com", Port: "10086", SpaceUsed: 0, SpaceFree: 1024 * 7, SpaceConsumeRate: 1024},
			&InfluxD{Host: "1one.pandora.com", Port: "11086", SpaceUsed: 0, SpaceFree: 1024 * 6, SpaceConsumeRate: 1024},
			&InfluxD{Host: "1one.pandora.com", Port: "12086", SpaceUsed: 0, SpaceFree: 1024 * 5, SpaceConsumeRate: 1024},
		},
		{
			&InfluxD{Host: "2two.pandora.com", Port: "08086", SpaceUsed: 0, SpaceFree: 1024 * 10, SpaceConsumeRate: 1024},
			&InfluxD{Host: "2two.pandora.com", Port: "09086", SpaceUsed: 0, SpaceFree: 1024 * 6, SpaceConsumeRate: 1024},
			&InfluxD{Host: "2two.pandora.com", Port: "10086", SpaceUsed: 0, SpaceFree: 1024 * 5, SpaceConsumeRate: 1024},
			&InfluxD{Host: "2two.pandora.com", Port: "11086", SpaceUsed: 0, SpaceFree: 1024 * 4, SpaceConsumeRate: 1024},
			&InfluxD{Host: "2two.pandora.com", Port: "12086", SpaceUsed: 0, SpaceFree: 1024 * 3, SpaceConsumeRate: 1024},
		},
		{
			&InfluxD{Host: "3three.pandora.com", Port: "08086", SpaceUsed: 0, SpaceFree: 1024 * 10, SpaceConsumeRate: 1024},
			&InfluxD{Host: "3three.pandora.com", Port: "09086", SpaceUsed: 0, SpaceFree: 1024 * 5, SpaceConsumeRate: 1024},
			&InfluxD{Host: "3three.pandora.com", Port: "10086", SpaceUsed: 0, SpaceFree: 1024 * 4, SpaceConsumeRate: 1024},
			&InfluxD{Host: "3three.pandora.com", Port: "11086", SpaceUsed: 0, SpaceFree: 1024 * 3, SpaceConsumeRate: 1024},
			&InfluxD{Host: "3three.pandora.com", Port: "12086", SpaceUsed: 0, SpaceFree: 1024 * 2, SpaceConsumeRate: 1024},
		},
		{
			&InfluxD{Host: "4four.pandora.com", Port: "08086", SpaceUsed: 0, SpaceFree: 1024 * 9, SpaceConsumeRate: 1024},
			&InfluxD{Host: "4four.pandora.com", Port: "09086", SpaceUsed: 0, SpaceFree: 1024 * 8, SpaceConsumeRate: 1024},
			&InfluxD{Host: "4four.pandora.com", Port: "10086", SpaceUsed: 0, SpaceFree: 1024 * 7, SpaceConsumeRate: 1024},
			&InfluxD{Host: "4four.pandora.com", Port: "11086", SpaceUsed: 0, SpaceFree: 1024 * 6, SpaceConsumeRate: 1024},
			&InfluxD{Host: "4four.pandora.com", Port: "12086", SpaceUsed: 0, SpaceFree: 1024 * 5, SpaceConsumeRate: 1024},
		},
		{
			&InfluxD{Host: "5five.pandora.com", Port: "08086", SpaceUsed: 0, SpaceFree: 1024 * 4, SpaceConsumeRate: 1024},
			&InfluxD{Host: "5five.pandora.com", Port: "09086", SpaceUsed: 0, SpaceFree: 1024 * 3, SpaceConsumeRate: 1024},
			&InfluxD{Host: "5five.pandora.com", Port: "10086", SpaceUsed: 0, SpaceFree: 1024 * 2, SpaceConsumeRate: 1024},
			&InfluxD{Host: "5five.pandora.com", Port: "11086", SpaceUsed: 0, SpaceFree: 1024 * 2, SpaceConsumeRate: 1024},
			&InfluxD{Host: "5five.pandora.com", Port: "12086", SpaceUsed: 0, SpaceFree: 1024 * 1, SpaceConsumeRate: 1024},
		},
	}

	influxds := make([]*InfluxD, 0)
	for _, in := range influxs {
		influxds = append(influxds, in...)
	}

	got, err := PickInfluxds(influxds, 3)
	if err != nil {
		t.Error(err)
	}
	exp := []*InfluxD{influxs[0][0], influxs[1][0], influxs[2][0]}
	if !reflect.DeepEqual(got, exp) {
		t.Errorf("make group fail\ngot:%v\nexp:%v\n", spew.Sdump(got), spew.Sdump(exp))
	}
	color(got)

	//第二次make group
	got, err = PickInfluxds(influxds, 3)
	if err != nil {
		t.Error(err)
	}
	exp = []*InfluxD{influxs[0][3], influxs[1][1], influxs[3][3]}

	if !reflect.DeepEqual(got, exp) {
		t.Errorf("make group fail\ngot:%v\nexp:%v\n", spew.Sdump(got), spew.Sdump(exp))
	}

	color(got)

	//第三次make group
	got, err = PickInfluxds(influxds, 3)
	if err != nil {
		t.Error(err)
	}
	exp = []*InfluxD{influxs[0][4], influxs[1][2], influxs[2][1]}

	if !reflect.DeepEqual(got, exp) {
		t.Errorf("make group fail\ngot:%v\nexp:%v\n", spew.Sdump(got), spew.Sdump(exp))
	}
}
