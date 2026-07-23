// magnifier_rs.rs — цифровая лупа с заморозкой кадра на Rust (opencv-rust)

extern crate opencv;
use opencv::{
    core::{Mat, Point, Scalar, Size, Vector},
    highgui::{imshow, wait_key, named_window, destroy_window, set_mouse_callback, MouseEvent, MouseEventFlags},
    imgproc::{resize, rectangle, circle, put_text, FONT_HERSHEY_SIMPLEX},
    videoio::{VideoCapture, CAP_ANY},
    prelude::*,
};

fn main() -> opencv::Result<()> {
    let mut cap = VideoCapture::new(0, CAP_ANY)?;
    if !cap.is_opened()? {
        println!("Не удалось открыть камеру");
        return Ok(());
    }

    named_window("Digital Magnifier", 0)?;
    let mut frozen = false;
    let mut frozen_frame = Mat::default();
    let mut magnifier_pos = Point::new(200, 200);
    let mut zoom = 2.0f64;
    let magnifier_size = 150;
    let mut frame = Mat::default();

    println!("🔍 Цифровая лупа (заморозка кадра)");
    println!("Пробел — заморозить/возобновить");
    println!("Колесо мыши — изменение зума (не реализовано)");
    println!("S — сохранить кадр");
    println!("Q — выход");

    // Обработка мыши (упрощённо)
    // В opencv-rust можно использовать set_mouse_callback, но для простоты оставим перемещение клавишами.

    loop {
        if !frozen {
            cap.read(&mut frame)?;
            if frame.empty() { break; }
        } else {
            frame = frozen_frame.clone();
        }

        let mut display = frame.clone();

        if !frozen {
            let size = (magnifier_size as f64 / zoom) as i32;
            let rect = opencv::core::Rect::new(
                magnifier_pos.x - size/2,
                magnifier_pos.y - size/2,
                size, size
            );
            rectangle(&mut display, rect, Scalar::new(0.0, 255.0, 0.0, 0.0), 2, 8, 0)?;
        } else {
            // Увеличение области
            let x = magnifier_pos.x;
            let y = magnifier_pos.y;
            let region_size = (magnifier_size as f64 / zoom) as i32;
            let x1 = x - region_size/2;
            let y1 = y - region_size/2;
            let x2 = x + region_size/2;
            let y2 = y + region_size/2;
            let (cols, rows) = (frame.cols(), frame.rows());
            let x1 = x1.max(0).min(cols);
            let y1 = y1.max(0).min(rows);
            let x2 = x2.max(0).min(cols);
            let y2 = y2.max(0).min(rows);
            if x2 > x1 && y2 > y1 {
                let region = frame.roi(opencv::core::Rect::new(x1, y1, x2-x1, y2-y1))?;
                let mut mag = Mat::default();
                resize(&region, &mut mag, Size::new(magnifier_size, magnifier_size), 0.0, 0.0, 1)?;
                // Показываем в отдельном окне
                imshow("Magnifier View", &mag)?;
                circle(&mut display, magnifier_pos, magnifier_size/2, Scalar::new(0.0, 255.0, 255.0, 0.0), 2, 8, 0)?;
            }
        }

        // Информация
        let status = if frozen { "FROZEN" } else { "LIVE" };
        let color = if frozen { Scalar::new(0.0, 0.0, 255.0, 0.0) } else { Scalar::new(0.0, 255.0, 0.0, 0.0) };
        put_text(&mut display, format!("Zoom: {:.1}", zoom).as_str(),
                 Point::new(10, 30), FONT_HERSHEY_SIMPLEX, 0.7,
                 Scalar::new(255.0, 255.0, 255.0, 0.0), 2, 8, false)?;
        put_text(&mut display, status, Point::new(10, 60),
                 FONT_HERSHEY_SIMPLEX, 0.7, color, 2, 8, false)?;

        imshow("Digital Magnifier", &display)?;

        let key = wait_key(1)?;
        if key == 'q' as i32 { break; }
        else if key == ' ' as i32 {
            if frozen {
                frozen = false;
                destroy_window("Magnifier View")?;
            } else {
                frozen_frame = frame.clone();
                frozen = true;
            }
        } else if key == 's' as i32 {
            let filename = if frozen { "frozen_frame.png" } else { "live_frame.png" };
            let _ = opencv::imgcodecs::imwrite(filename, if frozen { &frozen_frame } else { &frame }, &Vector::default())?;
            println!("Сохранено: {}", filename);
        } else if key == '+' as i32 {
            zoom *= 1.2;
            if zoom > 10.0 { zoom = 10.0; }
        } else if key == '-' as i32 {
            zoom /= 1.2;
            if zoom < 1.0 { zoom = 1.0; }
        } else if key == 119 { // w
            magnifier_pos.y -= 10;
        } else if key == 97 { // a
            magnifier_pos.x -= 10;
        } else if key == 115 { // s
            magnifier_pos.y += 10;
        } else if key == 100 { // d
            magnifier_pos.x += 10;
        }
        // Ограничение
        let (cols, rows) = (frame.cols(), frame.rows());
        if magnifier_pos.x < 0 { magnifier_pos.x = 0; }
        if magnifier_pos.y < 0 { magnifier_pos.y = 0; }
        if magnifier_pos.x > cols { magnifier_pos.x = cols; }
        if magnifier_pos.y > rows { magnifier_pos.y = rows; }
    }
    Ok(())
}
