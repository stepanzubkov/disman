package main

import (
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type CT = layout.Context
type D = layout.Dimensions

func main() {
	go func() {
		window := new(app.Window)
		err := run(window)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
    passwd := Getpwnam("stepan")
    log.Println(passwd.Shell)
	app.Main()
}


func run(window *app.Window) error {
	theme := material.NewTheme()
	var ops op.Ops
    var loginButton widget.Clickable
    var usernameInput widget.Editor
    usernameInput.SingleLine = true
    var passwordInput widget.Editor
    passwordInput.SingleLine = true
    passwordInput.Mask = '*'

    var username string
    var password string
	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

            if loginButton.Clicked(gtx) {
                username = usernameInput.Text()
                password = passwordInput.Text()
                log.Println("username: ", username)
                log.Println("password: ", password)
                usernameInput.SetText("")
                passwordInput.SetText("")
                Login(username, password)
                log.Fatalln("End")
            }

            flex := layout.Flex{
                Axis: layout.Vertical,
                Spacing: layout.SpaceSides,
            }

            flex.Layout(gtx,
                layout.Rigid(InputLayout(gtx, theme, &usernameInput, "Username")),
                layout.Rigid(InputLayout(gtx, theme, &passwordInput, "Password")),
                layout.Rigid(
                    func(gtx CT) D {
                        margins := layout.Inset{
                            Right: unit.Dp(200),
                            Left: unit.Dp(200),
                        }
                        return margins.Layout(gtx,
                            func(gtx CT) D {
                                btn := material.Button(theme, &loginButton, "Login")
                                return btn.Layout(gtx)
                            },
                        )
                    },
                ),
            )

			e.Frame(gtx.Ops)
		}
	}
}
func InputLayout(gtx CT, theme *material.Theme, input *widget.Editor, hint string) (func(gtx CT) D) {
    return func(gtx CT) D {
        margins := layout.Inset{
            Right: unit.Dp(200),
            Left: unit.Dp(200),
            Bottom: unit.Dp(20),
        }
        border := widget.Border{
              Color:        color.NRGBA{R: 204, G: 204, B: 204, A: 255},
              CornerRadius: unit.Dp(3),
              Width:        unit.Dp(2),
        }
        inset := layout.Inset{
                Top: 10,
                Right: 10,
                Bottom: 5,
                Left: 10,
        }
        return margins.Layout(gtx,
            func(gtx CT) D {
                ed := material.Editor(theme, input, hint)
                return border.Layout(gtx,
                    func(gtx CT) D {
                        return inset.Layout(gtx, ed.Layout)
                    },
                )

            },
        )
    }
}
