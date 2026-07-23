// magnifier_cpp.cpp — цифровая лупа с заморозкой кадра на C++ (OpenCV)

#include <opencv2/opencv.hpp>
#include <iostream>

using namespace cv;
using namespace std;

class DigitalMagnifier {
private:
    VideoCapture cap;
    Mat frame, frozenFrame;
    bool frozen;
    Point magnifierPos;
    double zoom;
    int magnifierSize;
    string windowName, magWindow;

public:
    DigitalMagnifier() : frozen(false), magnifierPos(200, 200), zoom(2.0), magnifierSize(150),
                         windowName("Digital Magnifier"), magWindow("Magnifier View") {
        cap.open(0);
        if (!cap.isOpened()) {
            cerr << "Не удалось открыть камеру" << endl;
            exit(1);
        }
        namedWindow(windowName);
        setMouseCallback(windowName, mouseCallback, this);
        cout << "🔍 Цифровая лупа (заморозка кадра)" << endl;
        cout << "Пробел — заморозить/возобновить" << endl;
        cout << "Колесо мыши — изменение зума" << endl;
        cout << "S — сохранить кадр" << endl;
        cout << "Q — выход" << endl;
    }

    static void mouseCallback(int event, int x, int y, int flags, void* userdata) {
        auto* app = static_cast<DigitalMagnifier*>(userdata);
        if (event == EVENT_MOUSEMOVE) {
            app->magnifierPos = Point(x, y);
        } else if (event == EVENT_MOUSEWHEEL) {
            int delta = getMouseWheelDelta(flags);
            app->zoom += delta / 120.0 * 0.5;
            app->zoom = max(1.0, min(10.0, app->zoom));
        }
    }

    Mat getMagnifiedRegion(const Mat& src) {
        int x = magnifierPos.x, y = magnifierPos.y;
        int regionSize = (int)(magnifierSize / zoom);
        int x1 = x - regionSize/2, y1 = y - regionSize/2;
        int x2 = x + regionSize/2, y2 = y + regionSize/2;
        // Проверка границ
        x1 = max(0, min(src.cols, x1));
        y1 = max(0, min(src.rows, y1));
        x2 = max(0, min(src.cols, x2));
        y2 = max(0, min(src.rows, y2));
        if (x2 <= x1 || y2 <= y1) return Mat();
        Mat region = src(Rect(x1, y1, x2-x1, y2-y1));
        Mat magnified;
        resize(region, magnified, Size(magnifierSize, magnifierSize), 0, 0, INTER_LINEAR);
        // Рисуем центр
        circle(magnified, Point(magnifierSize/2, magnifierSize/2), 5, Scalar(0, 255, 0), 1);
        return magnified;
    }

    void run() {
        while (true) {
            if (!frozen) {
                cap >> frame;
                if (frame.empty()) break;
            } else {
                frame = frozenFrame.clone();
            }

            Mat display = frame.clone();
            if (!frozen) {
                // Рисуем прямоугольник-указатель
                int size = (int)(magnifierSize / zoom);
                rectangle(display, Point(magnifierPos.x-size/2, magnifierPos.y-size/2),
                          Point(magnifierPos.x+size/2, magnifierPos.y+size/2), Scalar(0, 255, 0), 2);
            } else {
                // Показываем увеличенную область
                Mat mag = getMagnifiedRegion(frozenFrame);
                if (!mag.empty()) {
                    imshow(magWindow, mag);
                    // Рисуем круг на основном кадре
                    circle(display, magnifierPos, magnifierSize/2, Scalar(0, 255, 255), 2);
                }
            }

            // Информация
            string status = frozen ? "FROZEN" : "LIVE";
            putText(display, "Zoom: " + to_string(zoom), Point(10, 30),
                    FONT_HERSHEY_SIMPLEX, 0.7, Scalar(255,255,255), 2);
            putText(display, status, Point(10, 60),
                    FONT_HERSHEY_SIMPLEX, 0.7, frozen ? Scalar(0,0,255) : Scalar(0,255,0), 2);

            imshow(windowName, display);

            char key = waitKey(1) & 0xFF;
            if (key == 'q') break;
            else if (key == ' ') {
                if (frozen) {
                    frozen = false;
                    destroyWindow(magWindow);
                } else {
                    if (!frame.empty()) {
                        frozenFrame = frame.clone();
                        frozen = true;
                    }
                }
            } else if (key == 's') {
                string filename = frozen ? "frozen_frame.png" : "live_frame.png";
                imwrite(filename, frozen ? frozenFrame : frame);
                cout << "Сохранено: " << filename << endl;
            }
        }
        cap.release();
        destroyAllWindows();
    }
};

int main() {
    DigitalMagnifier app;
    app.run();
    return 0;
}
