package main

import (
	"fmt"

	imgui "github.com/AllenDang/cimgui-go"
)

func FilterItem(index int) {
	if imgui.Button(fmt.Sprintf("X (%d)", index)) {
		Filters = append(Filters[:index], Filters[index+1:]...)
		return
	}

	if imgui.IsItemHovered() {
		imgui.SetTooltip("delete filter")
	}

	imgui.SameLine()

	imgui.SetNextItemWidth(160)
	if imgui.BeginCombo(fmt.Sprintf("format%d", index), FilterNames[Filters[index].What]) {
		for i, val := range FilterNames {
			if (imgui.SelectableBoolV(val, false, 0, imgui.Vec2{X: 160, Y: 0})) {
				Filters[index].What = FilterType(i)
			}
		}

		imgui.EndCombo()
	}

	if Filters[index].What == FilterUpscale {
		imgui.SameLine()

		imgui.SetNextItemWidth(160)
		imgui.InputInt(fmt.Sprintf("upscale (%d)", index), &(Filters[index].IntFactor))
	}
}
