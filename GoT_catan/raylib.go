package main

import (
	"fmt"
	"math"
    "strings"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	RESET = iota
	DRAW_1
	DRAW_2
	UNDO
	REDO
)

type Button struct {
	x      int32
	y      int32
	width  int32
	height int32
	text   string
	action int
}

type Action struct {
	action int
	state [9]int
	tokenCount int
	tokensDrawn int
	foundTokenNames []int
}


func drawButtons(buttons []Button) {
	for _, button := range buttons {
		rl.DrawRectangle(button.x+1, button.y+1, button.width-3, button.height-3, rl.LightGray)

		rl.DrawText(button.text, button.x+10, button.y+5, button.height-10, rl.Black)
	}
}

func getTokenName(token int) string {
	if token == -1 {
		return "Nothing"
	}
	var name string
	switch token / 3 {
	case 0:
		name += "Foot"
	case 1:
		name += "Caves"
	case 2:
		name += "Ladder"
	}
	name += " "
	switch token % 3 {
	case 0:
		name += "infantry"
	case 1:
		name += "climber"
	case 2:
		name += "giant"
	}
	return name
}

func takeToken(wildingTokens *[9]int, tokenCount *int) int {
	if tokenCount == nil || *tokenCount == 0 {
		return -1
	}
	value := rl.GetRandomValue(0, int32(*tokenCount-1))

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

func textCentered(text string, x, y, size int32, color rl.Color) {
	rl.DrawText(text, x-int32(int(size)*len(text)/4), y, size, color)
}

func getTexture(token int, textures [6]rl.Texture2D) (rl.Texture2D, rl.Texture2D) {
	var texture, texture_type rl.Texture2D
	switch token / 3 {
	case 0:
		texture = textures[0] //foot_texture
	case 1:
		texture = textures[1] //cave_texture
	case 2:
		texture = textures[2] //ladder_texture
	}

	switch token % 3 {
	case 0:
		texture_type = textures[3] //regular_texture
	case 1:
		texture_type = textures[4] //climber_texture
	case 2:
		texture_type = textures[5] //giant_texture
	}
	return texture, texture_type
}

func drawBigToken(centerX, centerY int32, token int, textures [6]rl.Texture2D) {
	var texture_camp rl.Texture2D
	var texture_type rl.Texture2D
	texture_camp, texture_type = getTexture(token, textures)
	var scale float32 = 3
	var offset int32 = -60
	rl.DrawRectangleRounded(rl.Rectangle{X: float32(centerX - int32(scale*12.5)), Y: float32(centerY - int32(scale*12.5) + offset), Width: 25 * scale, Height: 25 * scale}, 0.5, 0, rl.Gray)
	rl.DrawRectangleRounded(rl.Rectangle{X: float32(centerX - int32(scale*11)), Y: float32(centerY - int32(scale*11) + offset), Width: 22 * scale, Height: 22 * scale}, 0.5, 0, rl.LightGray)

	rl.DrawTextureEx(texture_camp, rl.Vector2{X: float32(centerX - int32(scale*10)), Y: float32(centerY - int32(scale*10) + offset)}, 0, scale, rl.White)
	rl.DrawTextureEx(texture_type, rl.Vector2{X: float32(centerX + int32(15*scale) - int32(scale*10)), Y: float32(centerY - int32(scale*10) + offset)}, 0, scale/2, rl.White)
}

func main() {

	var screenWidth int32 = 800
	var screenHeight int32 = 450
	wildingTokens := [9]int{12, 4, 3, 12, 4, 3, 12, 4, 3}
	tokenCount := 57
	lastDraw := 0

	foundTokenNames := []int{}
	undoList := []Action{}
	redoList := []Action{}

	rl.InitWindow(screenWidth, screenHeight, "raylib [core] example - basic window")
	defer rl.CloseWindow()

	foot_image := rl.LoadImage("foot.png")
	ladder_image := rl.LoadImage("ladder.png")
	cave_image := rl.LoadImage("cave.png")

	foot_texture := rl.LoadTextureFromImage(foot_image)
	ladder_texture := rl.LoadTextureFromImage(ladder_image)
	cave_texture := rl.LoadTextureFromImage(cave_image)

	rl.UnloadImage(foot_image)
	rl.UnloadImage(ladder_image)
	rl.UnloadImage(cave_image)

	regular_image := rl.LoadImage("regular.png")
	climber_image := rl.LoadImage("climber.png")
	giant_image := rl.LoadImage("giant.png")

	regular_texture := rl.LoadTextureFromImage(regular_image)
	climber_texture := rl.LoadTextureFromImage(climber_image)
	giant_texture := rl.LoadTextureFromImage(giant_image)

	rl.UnloadImage(regular_image)
	rl.UnloadImage(climber_image)
	rl.UnloadImage(giant_image)

	rl.SetTargetFPS(60)
	buttons := []Button{}
	buttons = append(buttons, Button{x: screenWidth - 100, y: screenHeight - 50, width: 80, height: 30, text: "Reset", action: RESET},
		Button{x: 20, y: screenHeight - 50, width: 100, height: 30, text: "Draw 1", action: DRAW_1},
		Button{x: 130, y: screenHeight - 50, width: 100, height: 30, text: "Draw 2", action: DRAW_2},
		Button{x: screenWidth - 210, y: screenHeight - 44, width: 50, height: 24, text: "undo", action: UNDO},
		Button{x: screenWidth - 160, y: screenHeight - 44, width: 50, height: 24, text: "redo", action: REDO})

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		//GUI
		textCentered("Game of thrones token generator", screenWidth/2-10, 80, 20, rl.LightGray)
		x_pos := int32(375 * (1 + math.Sin(rl.GetTime())))
		rl.DrawRectangleGradientV(x_pos, 240, 50, 50, rl.Red, rl.Blue)

		textCentered(strings.Trim(strings.Replace(fmt.Sprint(wildingTokens), " ", ",", -1), "[]"), screenWidth/2+5, 350, 16, rl.LightGray)
		textCentered(fmt.Sprint(tokenCount)+ " left", screenWidth/2, 370, 16, rl.LightGray)

		drawButtons(buttons)

		if rl.IsMouseButtonReleased(rl.MouseLeftButton) {
			for i := len(buttons) - 1; i >= 0; i-- {
				button := buttons[i]
				if rl.GetMousePosition().X >= float32(button.x) && rl.GetMousePosition().X <= float32(button.x+button.width) && rl.GetMousePosition().Y >= float32(button.y) && rl.GetMousePosition().Y <= float32(button.y+button.height) {
					//Save state
					
					undoLength := 100
					if len(undoList) > undoLength {
						undoList = undoList[len(undoList)-undoLength:]
					}

					if button.action != UNDO && button.action != REDO {
						redoList = []Action{}
						action := Action{action: button.action, state: wildingTokens, tokenCount: tokenCount, tokensDrawn: lastDraw, foundTokenNames: foundTokenNames}
						undoList = append(undoList, action)
					}

					fmt.Println("Pressed ", button.text)

					switch button.action {
					case RESET:
						fmt.Println("Tokens reset!")
						wildingTokens = [9]int{12, 4, 3, 12, 4, 3, 12, 4, 3}
						tokenCount = 57
						foundTokenNames = []int{}
						lastDraw = 0

					case DRAW_1:
						token := takeToken(&wildingTokens, &tokenCount)
						foundTokenNames = append(foundTokenNames, token)
						fmt.Println("Drew:", getTokenName(token))
						lastDraw = 1

					case DRAW_2:
						token := takeToken(&wildingTokens, &tokenCount)
						foundTokenNames = append(foundTokenNames, token)
						fmt.Println("Drew:", getTokenName(token))
						token = takeToken(&wildingTokens, &tokenCount)
						foundTokenNames = append(foundTokenNames, token)
						fmt.Println("Drew:", getTokenName(token))
						lastDraw = 2
					case UNDO:
						if len(undoList) > 0 {
							redoList = append(redoList, Action{action: button.action, state: wildingTokens, tokenCount: tokenCount, tokensDrawn: lastDraw, foundTokenNames: foundTokenNames})
							wildingTokens = undoList[len(undoList)-1].state
							tokenCount = undoList[len(undoList)-1].tokenCount
							foundTokenNames = undoList[len(undoList)-1].foundTokenNames
							lastDraw = undoList[len(undoList)-1].tokensDrawn
							undoList = undoList[:len(undoList)-1]						
						}
					case REDO:
						if len(redoList) > 0 {
							undoList = append(undoList, Action{action: button.action, state: wildingTokens, tokenCount: tokenCount, tokensDrawn: lastDraw, foundTokenNames: foundTokenNames})
							wildingTokens = redoList[len(redoList)-1].state
							tokenCount = redoList[len(redoList)-1].tokenCount
							foundTokenNames = redoList[len(redoList)-1].foundTokenNames
							lastDraw = redoList[len(redoList)-1].tokensDrawn
							redoList = redoList[:len(redoList)-1]	
						}
					}
					break
				}

			}
		}

		// Lightblue box around last draws
		switch lastDraw {
		case 1:
			rl.DrawRectangle(10, 10, 45, 20, rl.SkyBlue)
		case 2:
			rl.DrawRectangle(10, 10, 45, 40, rl.SkyBlue)
		}

		// Draw last 20 tokens
		for i := len(foundTokenNames) - 1; i >= max(0, len(foundTokenNames)-19); i-- {
			token := foundTokenNames[i]

			texture_camp, texture_type := getTexture(token, [6]rl.Texture2D{foot_texture, cave_texture, ladder_texture, regular_texture, climber_texture, giant_texture})

			height := int32(len(foundTokenNames) - i - 1)
			rl.DrawTexture(texture_camp, 10, 10+height*20, rl.ColorAlpha(rl.White, 1.1-(float32(height-13)/5.0)))
			rl.DrawTexture(texture_type, 35, 10+height*20, rl.ColorAlpha(rl.White, 1.1-(float32(height-13)/5.0)))

		}

		// Draw the tokens largely in the center
		if lastDraw == 1 {
			drawBigToken(screenWidth/2, screenHeight/2, foundTokenNames[len(foundTokenNames)-1], [6]rl.Texture2D{foot_texture, cave_texture, ladder_texture, regular_texture, climber_texture, giant_texture})
		}
		if lastDraw == 2 {
			drawBigToken(screenWidth/2-80, screenHeight/2, foundTokenNames[len(foundTokenNames)-1], [6]rl.Texture2D{foot_texture, cave_texture, ladder_texture, regular_texture, climber_texture, giant_texture})
			drawBigToken(screenWidth/2+80, screenHeight/2, foundTokenNames[len(foundTokenNames)-2], [6]rl.Texture2D{foot_texture, cave_texture, ladder_texture, regular_texture, climber_texture, giant_texture})
		}


		rl.EndDrawing()
	}
}
