package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/theme"
	"wyx0k.com/img-shadow/shadow"
)

var font *fyne.StaticResource

type imgShadowTheme struct{}

func init() {
	fontBytes, err := shadow.Asset("SiYuanHeiTiJiuZiXing-Regular-2.ttf")
	if err != nil {
		return
	}
	font = &fyne.StaticResource{
		StaticName:    "SiYuanHeiTiJiuZiXing-Regular-2.ttf",
		StaticContent: fontBytes,
	}
	fmt.Print("asd")
	fmt.Printf("%d", len(font.StaticContent))
}

// return bundled font resource
func (imgShadowTheme) TextFont() fyne.Resource { return font }
func (imgShadowTheme) TextBoldFont() fyne.Resource {
	return font
}

func (imgShadowTheme) BackgroundColor() color.Color { return theme.LightTheme().BackgroundColor() }
func (imgShadowTheme) ButtonColor() color.Color     { return theme.LightTheme().ButtonColor() }
func (imgShadowTheme) DisabledButtonColor() color.Color {
	return theme.LightTheme().DisabledButtonColor()
}
func (imgShadowTheme) IconColor() color.Color         { return theme.LightTheme().IconColor() }
func (imgShadowTheme) DisabledIconColor() color.Color { return theme.LightTheme().DisabledIconColor() }
func (imgShadowTheme) HyperlinkColor() color.Color    { return theme.LightTheme().HyperlinkColor() }
func (imgShadowTheme) TextColor() color.Color         { return theme.LightTheme().TextColor() }
func (imgShadowTheme) DisabledTextColor() color.Color { return theme.LightTheme().DisabledTextColor() }
func (imgShadowTheme) HoverColor() color.Color        { return theme.LightTheme().HoverColor() }
func (imgShadowTheme) PlaceHolderColor() color.Color  { return theme.LightTheme().PlaceHolderColor() }
func (imgShadowTheme) PrimaryColor() color.Color      { return theme.LightTheme().PrimaryColor() }
func (imgShadowTheme) FocusColor() color.Color        { return theme.LightTheme().FocusColor() }
func (imgShadowTheme) ScrollBarColor() color.Color    { return theme.LightTheme().ScrollBarColor() }
func (imgShadowTheme) ShadowColor() color.Color       { return theme.LightTheme().ShadowColor() }
func (imgShadowTheme) TextSize() int                  { return theme.LightTheme().TextSize() }
func (imgShadowTheme) TextItalicFont() fyne.Resource  { return theme.LightTheme().TextItalicFont() }
func (imgShadowTheme) TextBoldItalicFont() fyne.Resource {
	return theme.LightTheme().TextBoldItalicFont()
}
func (imgShadowTheme) TextMonospaceFont() fyne.Resource { return theme.LightTheme().TextMonospaceFont() }
func (imgShadowTheme) Padding() int                     { return theme.LightTheme().Padding() }
func (imgShadowTheme) IconInlineSize() int              { return theme.LightTheme().IconInlineSize() }
func (imgShadowTheme) ScrollBarSize() int               { return theme.LightTheme().ScrollBarSize() }
func (imgShadowTheme) ScrollBarSmallSize() int          { return theme.LightTheme().ScrollBarSmallSize() }
