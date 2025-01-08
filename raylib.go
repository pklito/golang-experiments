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
	action int
}

func drawGUI() {
	rl.DrawText("Congrats! You created your first window!", 190, 200, 20, rl.LightGray)
	rl.DrawText("Line two", 190, 220, 20, rl.LightGray)
	x_pos := int32(375 * (1 + math.Sin(rl.GetTime())))
	rl.DrawRectangleGradientV(x_pos, 200, 50, 50, rl.Red, rl.Blue)
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

	rl.DrawTextureEx(texture_camp, rl.Vector2{X: float32(centerX-int32(scale*10)), Y: float32(centerY-int32(scale*10)+offset)}, 0, scale, rl.White)
	rl.DrawTextureEx(texture_type, rl.Vector2{X: float32(centerX+int32(15*scale)-int32(scale*10)), Y: float32(centerY-int32(scale*10)+offset)}, 0, scale/2, rl.White)
}

func main() {

	var screenWidth int32 = 800
	var screenHeight int32 = 450
	wildingTokens := [9]int{8, 4, 4, 8, 4, 4, 8, 4, 4}
	tokenCount := 48
	lastDraw := 0
	lastDrawTime := -1000.0

	foundTokenNames := []int{}

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
		Button{x: 130, y: screenHeight - 50, width: 100, height: 30, text: "Draw 2", action: DRAW_2})

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		drawGUI()
		drawButtons(buttons)

		if rl.IsMouseButtonReleased(rl.MouseLeftButton) {
			for i := len(buttons) - 1; i >= 0; i-- {
				button := buttons[i]
				if rl.GetMousePosition().X >= float32(button.x) && rl.GetMousePosition().X <= float32(button.x+button.width) && rl.GetMousePosition().Y >= float32(button.y) && rl.GetMousePosition().Y <= float32(button.y+button.height) {
					fmt.Println("Pressed ", button.text)

					switch button.action {
					case RESET:
						fmt.Println("Tokens reset!")
						wildingTokens = [9]int{8, 4, 4, 8, 4, 4, 8, 4, 4}
						tokenCount = 48
						foundTokenNames = []int{}
						lastDraw = 0
						lastDrawTime = -1000.0

					case DRAW_1:
						token := takeToken(&wildingTokens, &tokenCount)
						foundTokenNames = append(foundTokenNames, token)
						fmt.Println("Drew:", getTokenName(token))
						lastDraw = 1
						lastDrawTime = rl.GetTime()

					case DRAW_2:
						token := takeToken(&wildingTokens, &tokenCount)
						foundTokenNames = append(foundTokenNames, token)
						fmt.Println("Drew:", getTokenName(token))
						token = takeToken(&wildingTokens, &tokenCount)
						foundTokenNames = append(foundTokenNames, token)
						fmt.Println("Drew:", getTokenName(token))
						lastDraw = 2
						lastDrawTime = rl.GetTime()
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
		if lastDrawTime+100000 > rl.GetTime() {
			if lastDraw == 1 {
				drawBigToken(screenWidth/2, screenHeight/2, foundTokenNames[len(foundTokenNames)-1], [6]rl.Texture2D{foot_texture, cave_texture, ladder_texture, regular_texture, climber_texture, giant_texture})
			}
			if lastDraw == 2 {
				drawBigToken(screenWidth/2-80, screenHeight/2, foundTokenNames[len(foundTokenNames)-1], [6]rl.Texture2D{foot_texture, cave_texture, ladder_texture, regular_texture, climber_texture, giant_texture})
				drawBigToken(screenWidth/2+80, screenHeight/2, foundTokenNames[len(foundTokenNames)-2], [6]rl.Texture2D{foot_texture, cave_texture, ladder_texture, regular_texture, climber_texture, giant_texture})
			}
		}

		rl.EndDrawing()
	}
}
