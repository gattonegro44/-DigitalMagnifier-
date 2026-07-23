# magnifier_python.py — цифровая лупа с заморозкой кадра на Python (OpenCV)

import cv2
import numpy as np
import sys

class DigitalMagnifier:
    def __init__(self):
        self.cap = cv2.VideoCapture(0)
        if not self.cap.isOpened():
            print("Не удалось открыть камеру")
            sys.exit(1)
        self.frozen = False
        self.frozen_frame = None
        self.magnifier_pos = (200, 200)  # центр лупы (x, y)
        self.zoom = 2.0
        self.magnifier_size = 150  # размер окна лупы (пикселей)
        self.dragging = False
        self.window_name = "Digital Magnifier"
        self.magnifier_window = "Magnifier View"

        cv2.namedWindow(self.window_name)
        cv2.setMouseCallback(self.window_name, self.mouse_callback)

        print("🔍 Цифровая лупа (заморозка кадра)")
        print("Пробел — заморозить/возобновить")
        print("Колесо мыши — изменение зума")
        print("S — сохранить кадр")
        print("Q — выход")

    def mouse_callback(self, event, x, y, flags, param):
        if event == cv2.EVENT_MOUSEMOVE:
            self.magnifier_pos = (x, y)
        elif event == cv2.EVENT_MOUSEWHEEL:
            # Изменение зума
            delta = flags >> 16
            self.zoom += delta / 120 * 0.5
            self.zoom = max(1.0, min(10.0, self.zoom))
        elif event == cv2.EVENT_LBUTTONDOWN:
            self.dragging = True
        elif event == cv2.EVENT_LBUTTONUP:
            self.dragging = False

    def get_magnified_region(self, frame):
        x, y = self.magnifier_pos
        h, w = frame.shape[:2]
        # Размер области для увеличения (в исходном разрешении)
        region_size = int(self.magnifier_size / self.zoom)
        # Координаты верхнего левого угла области
        x1 = int(x - region_size / 2)
        y1 = int(y - region_size / 2)
        x2 = int(x + region_size / 2)
        y2 = int(y + region_size / 2)
        # Проверка границ
        x1 = max(0, min(w, x1))
        y1 = max(0, min(h, y1))
        x2 = max(0, min(w, x2))
        y2 = max(0, min(h, y2))
        if x2 <= x1 or y2 <= y1:
            return None
        region = frame[y1:y2, x1:x2]
        if region.size == 0:
            return None
        # Увеличение
        magnified = cv2.resize(region, (self.magnifier_size, self.magnifier_size),
                               interpolation=cv2.INTER_LINEAR)
        # Рамка и указатель центра
        cv2.circle(magnified, (self.magnifier_size//2, self.magnifier_size//2), 5, (0, 255, 0), 1)
        return magnified

    def run(self):
        while True:
            if not self.frozen:
                ret, frame = self.cap.read()
                if not ret:
                    break
                display = frame.copy()
            else:
                display = self.frozen_frame.copy()

            # Отображение лупы на основном кадре
            if not self.frozen:
                # Показываем квадрат-указатель на основном кадре
                x, y = self.magnifier_pos
                size = int(self.magnifier_size / self.zoom)
                cv2.rectangle(display, (x-size//2, y-size//2), (x+size//2, y+size//2), (0, 255, 0), 2)
            else:
                # Рисуем лупу на замороженном кадре
                magnified = self.get_magnified_region(self.frozen_frame)
                if magnified is not None:
                    # Показываем увеличенную область в отдельном окне
                    cv2.imshow(self.magnifier_window, magnified)
                    # Также рисуем круг-лупу на основном кадре
                    x, y = self.magnifier_pos
                    cv2.circle(display, (x, y), self.magnifier_size//2, (0, 255, 255), 2)

            # Информация на кадре
            cv2.putText(display, f"Zoom: {self.zoom:.1f}x", (10, 30),
                        cv2.FONT_HERSHEY_SIMPLEX, 0.7, (255, 255, 255), 2)
            status = "FROZEN" if self.frozen else "LIVE"
            cv2.putText(display, status, (10, 60),
                        cv2.FONT_HERSHEY_SIMPLEX, 0.7, (0, 0, 255) if self.frozen else (0, 255, 0), 2)

            cv2.imshow(self.window_name, display)

            key = cv2.waitKey(1) & 0xFF
            if key == ord('q'):
                break
            elif key == ord(' '):  # пробел
                if self.frozen:
                    self.frozen = False
                    cv2.destroyWindow(self.magnifier_window)
                else:
                    ret, frame = self.cap.read()
                    if ret:
                        self.frozen_frame = frame.copy()
                        self.frozen = True
            elif key == ord('s'):
                if self.frozen:
                    filename = "frozen_frame.png"
                    cv2.imwrite(filename, self.frozen_frame)
                    print(f"Кадр сохранён как {filename}")
                else:
                    ret, frame = self.cap.read()
                    if ret:
                        filename = "live_frame.png"
                        cv2.imwrite(filename, frame)
                        print(f"Кадр сохранён как {filename}")

        self.cap.release()
        cv2.destroyAllWindows()

if __name__ == "__main__":
    app = DigitalMagnifier()
    app.run()
