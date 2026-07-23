// magnifier_go.go — цифровая лупа с заморозкой кадра на Go (GoCV)

package main

import (
    "fmt"
    "gocv.io/x/gocv"
    "image"
    "os"
    "strconv"
)

func main() {
    webcam, err := gocv.OpenVideoCapture(0)
    if err != nil {
        fmt.Printf("Ошибка открытия камеры: %v\n", err)
        return
    }
    defer webcam.Close()

    window := gocv.NewWindow("Digital Magnifier")
    magWindow := gocv.NewWindow("Magnifier View")
    defer window.Close()
    defer magWindow.Close()

    img := gocv.NewMat()
    frozen := gocv.NewMat()
    frozenFrame := false
    var magnifierPos image.Point = image.Point{200, 200}
    zoom := 2.0
    magnifierSize := 150

    fmt.Println("🔍 Цифровая лупа (заморозка кадра)")
    fmt.Println("Пробел — заморозить/возобновить")
    fmt.Println("Колесо мыши — изменение зума (не реализовано в GoCV, используйте +/-)")
    fmt.Println("S — сохранить кадр")
    fmt.Println("Q — выход")

    // Обработка мыши для перемещения лупы (требуется GoCV с mouse callback, в упрощённой версии нет)
    // Здесь используется простой вариант: перемещение клавишами стрелок.

    for {
        if !frozenFrame {
            if webcam.Read(&img) {
                // ok
            } else {
                break
            }
        } else {
            img = frozen.Clone()
        }

        display := img.Clone()

        if !frozenFrame {
            // Рисуем прямоугольник
            size := int(float64(magnifierSize) / zoom)
            rect := image.Rect(magnifierPos.X-size/2, magnifierPos.Y-size/2,
                               magnifierPos.X+size/2, magnifierPos.Y+size/2)
            gocv.Rectangle(&display, rect, gocv.NewScalar(0, 255, 0, 0), 2)
        } else {
            // Увеличение области
            x := magnifierPos.X
            y := magnifierPos.Y
            regionSize := int(float64(magnifierSize) / zoom)
            x1 := x - regionSize/2
            y1 := y - regionSize/2
            x2 := x + regionSize/2
            y2 := y + regionSize/2
            // Проверка границ
            rows := frozen.Rows()
            cols := frozen.Cols()
            if x1 < 0 { x1 = 0 }
            if y1 < 0 { y1 = 0 }
            if x2 > cols { x2 = cols }
            if y2 > rows { y2 = rows }
            if x2 > x1 && y2 > y1 {
                region := frozen.Region(image.Rect(x1, y1, x2, y2))
                mag := gocv.NewMat()
                gocv.Resize(region, &mag, image.Point{magnifierSize, magnifierSize}, 0, 0, gocv.InterpolationLinear)
                magWindow.IMShow(mag)
                // Рисуем круг на основном кадре
                gocv.Circle(&display, magnifierPos, magnifierSize/2, gocv.NewScalar(0, 255, 255, 0), 2)
                mag.Close()
                region.Close()
            }
        }

        // Информация
        status := "LIVE"
        if frozenFrame {
            status = "FROZEN"
        }
        gocv.PutText(&display, "Zoom: "+strconv.FormatFloat(zoom, 'f', 1, 64), image.Point{10, 30},
                     gocv.FontHersheySimplex, 0.7, gocv.NewScalar(255,255,255,0), 2)
        color := gocv.NewScalar(0,255,0,0)
        if frozenFrame {
            color = gocv.NewScalar(0,0,255,0)
        }
        gocv.PutText(&display, status, image.Point{10, 60}, gocv.FontHersheySimplex, 0.7, color, 2)

        window.IMShow(display)

        key := window.WaitKey(1)
        if key == 113 { // 'q'
            break
        } else if key == 32 { // space
            if frozenFrame {
                frozenFrame = false
                magWindow.Close()
                magWindow = gocv.NewWindow("Magnifier View")
            } else {
                frozen = img.Clone()
                frozenFrame = true
            }
        } else if key == 115 { // 's'
            filename := "frozen_frame.png"
            if !frozenFrame {
                filename = "live_frame.png"
            }
            if ok := gocv.IMWrite(filename, frozenFrame ? frozen : img); ok {
                fmt.Printf("Сохранено: %s\n", filename)
            }
        } else if key == 43 { // '+'
            zoom *= 1.2
            if zoom > 10.0 { zoom = 10.0 }
        } else if key == 45 { // '-'
            zoom /= 1.2
            if zoom < 1.0 { zoom = 1.0 }
        } else if key == 0 { // arrow keys (not reliable)
            // Для перемещения используем стрелки, но в GoCV они не обрабатываются как обычные клавиши.
            // Можно использовать другие клавиши: 'w','a','s','d'
        }
        // Перемещение лупы клавишами w,a,s,d
        if key == 119 { // w
            magnifierPos.Y -= 10
        } else if key == 97 { // a
            magnifierPos.X -= 10
        } else if key == 115 { // s (down)
            magnifierPos.Y += 10
        } else if key == 100 { // d
            magnifierPos.X += 10
        }
        // Ограничиваем позицию в пределах кадра
        if magnifierPos.X < 0 { magnifierPos.X = 0 }
        if magnifierPos.Y < 0 { magnifierPos.Y = 0 }
        if magnifierPos.X > img.Cols() { magnifierPos.X = img.Cols() }
        if magnifierPos.Y > img.Rows() { magnifierPos.Y = img.Rows() }

        display.Close()
        img.Close()
    }
}
