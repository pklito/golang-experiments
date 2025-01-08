package main

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	RESET = iota
	DRAW_1
	DRAW_2

)

type Button struct {
	x      int32
	y      int32
	width  int32
	height int32
	text   string
	action	 int
}

func drawGUI(tokenNames []string) {
	rl.DrawText("Congrats! You created your first window!", 190, 200, 20, rl.LightGray)
	rl.DrawText("Line two", 190, 220, 20, rl.LightGray)
	x_pos := int32(375 * (1 + math.Sin(rl.GetTime())))
	rl.DrawRectangleGradientV(x_pos, 200, 50, 50, rl.Red, rl.Blue)

	for i, tokenName := range tokenNames {
		rl.DrawText(tokenName, 20, 20 + int32(i * 20), 20, rl.Black)
	}
	// test_image := rl.LoadImage("test.png")
	// texture := rl.LoadTextureFromImage(test_image);
	// rl.DrawTexture(texture, 50, 50, rl.White)
}

func drawButtons(buttons []Button) {
	for _, button := range buttons {
		 rl.DrawRectangle(button.x + 1, button.y + 1, button.width - 3, button.height - 3, rl.LightGray)

		rl.DrawText(button.text, button.x+10, button.y+5, button.height-10, rl.Black)
	}
}

func getTokenName(token int) string {
	if token == -1 {
		return "Nothing"
	}
	var name string
	switch token/3 {
	case 0:
		name += "Foot"
	case 1:
		name += "Caves"
	case 2:
		name += "Ladder"
	}
	name += " "
	switch token%3 {
	case 0:
		name += "infantry"
	case 1:
		name += "climber"
	case 2:
		name += "giant"
	}
	return name
}

func takeToken(wildingTokens *[9]int, tokenCount *int) int{
	if tokenCount == nil || *tokenCount == 0 {
		return -1
	}
	value := rl.GetRandomValue(0, int32(*tokenCount - 1))
	
	for i := 0; i < 9; i++ {
		value -= int32(wildingTokens[i])
		if value < 0 {
			wildingTokens[i] -= 1
			*tokenCount -= 1
			return i
		}
	}
	return -1
}

func main() {
	var screenWidth int32 = 800
	var screenHeight int32 = 450
	wildingTokens := [9]int{8,4,4,8,4,4,8,4,4}
	tokenCount := 48

	foundTokenNames := []string{}

	rl.InitWindow(screenWidth, screenHeight, "raylib [core] example - basic window")
	defer rl.CloseWindow()
 
	rl.SetTargetFPS(60)
	buttons := []Button{}
	buttons = append(buttons, Button{x: screenWidth - 100,y: screenHeight - 50, width: 80,height: 30,text: "Reset", action: RESET}, 
						Button{x: 20,y: screenHeight - 50,width: 100,height: 30,text: "Draw 1", action: DRAW_1}, 
						Button{x: 130,y: screenHeight - 50,width: 100,height: 30,text: "Draw 2", action: DRAW_2})
	
	
	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		drawGUI(foundTokenNames)
		drawButtons(buttons)

		if rl.IsMouseButtonReleased(rl.MouseLeftButton) {
			for i := len(buttons) - 1; i >= 0; i-- {
				button := buttons[i]
				if rl.GetMousePosition().X >= float32(button.x) && rl.GetMousePosition().X <= float32(button.x+button.width) && rl.GetMousePosition().Y >= float32(button.y) && rl.GetMousePosition().Y <= float32(button.y+button.height) {
					fmt.Println("Pressed ", button.text)
	
					switch button.action {
					case RESET:
						fmt.Println("Tokens reset!")
						wildingTokens = [9]int{8,4,4,8,4,4,8,4,4}
						tokenCount = 48
						
					case DRAW_1:
						token := takeToken(&wildingTokens, &tokenCount)
						foundTokenNames = append(foundTokenNames, getTokenName(token))
						fmt.Println("Drew:", token)

					case DRAW_2:
						token := takeToken(&wildingTokens, &tokenCount)
						foundTokenNames = append(foundTokenNames, getTokenName(token))
						fmt.Println("Drew:", token)
						token = takeToken(&wildingTokens, &tokenCount)
						foundTokenNames = append(foundTokenNames, getTokenName(token))
						fmt.Println("Drew:", token)
					}
					break
				}
	
			}
		}

		rl.EndDrawing()
	}
}
