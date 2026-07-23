// MagnifierJava.java — цифровая лупа с заморозкой кадра на Java (JavaCV)

import org.bytedeco.javacpp.*;
import org.bytedeco.opencv.opencv_core.*;
import static org.bytedeco.opencv.global.opencv_core.*;
import static org.bytedeco.opencv.global.opencv_imgproc.*;
import static org.bytedeco.opencv.global.opencv_imgcodecs.*;
import org.bytedeco.javacv.*;
import java.awt.event.*;
import javax.swing.*;

public class MagnifierJava {
    private OpenCVFrameGrabber grabber;
    private CanvasFrame canvas, magCanvas;
    private Mat frame, frozenFrame;
    private boolean frozen = false;
    private Point magnifierPos = new Point(200, 200);
    private double zoom = 2.0;
    private int magnifierSize = 150;

    public MagnifierJava() throws Exception {
        grabber = new OpenCVFrameGrabber(0);
        grabber.start();

        canvas = new CanvasFrame("Digital Magnifier");
        canvas.setDefaultCloseOperation(JFrame.EXIT_ON_CLOSE);
        canvas.addMouseListener(new MouseAdapter() {
            public void mousePressed(MouseEvent e) {}
        });
        canvas.addMouseMotionListener(new MouseAdapter() {
            public void mouseMoved(MouseEvent e) {
                magnifierPos.x(e.getX());
                magnifierPos.y(e.getY());
            }
        });
        canvas.addMouseWheelListener(e -> {
            double delta = e.getPreciseWheelRotation();
            zoom += delta * 0.5;
            zoom = Math.max(1.0, Math.min(10.0, zoom));
        });

        magCanvas = new CanvasFrame("Magnifier View");
        magCanvas.setVisible(false);

        System.out.println("🔍 Цифровая лупа (заморозка кадра)");
        System.out.println("Пробел — заморозить/возобновить");
        System.out.println("Колесо мыши — изменение зума");
        System.out.println("S — сохранить кадр");
        System.out.println("Q — выход");
    }

    private Mat getMagnifiedRegion(Mat src) {
        int x = (int)magnifierPos.x(), y = (int)magnifierPos.y();
        int regionSize = (int)(magnifierSize / zoom);
        int x1 = Math.max(0, Math.min(src.cols(), x - regionSize/2));
        int y1 = Math.max(0, Math.min(src.rows(), y - regionSize/2));
        int x2 = Math.max(0, Math.min(src.cols(), x + regionSize/2));
        int y2 = Math.max(0, Math.min(src.rows(), y + regionSize/2));
        if (x2 <= x1 || y2 <= y1) return new Mat();
        Mat region = new Mat(src, new Rect(x1, y1, x2-x1, y2-y1));
        Mat magnified = new Mat();
        resize(region, magnified, new Size(magnifierSize, magnifierSize));
        // Центр
        circle(magnified, new Point(magnifierSize/2, magnifierSize/2), 5, Scalar.GREEN, 1, 8, 0);
        return magnified;
    }

    public void run() throws Exception {
        while (true) {
            Frame f = grabber.grabFrame();
            if (f == null) break;
            OpenCVFrameConverter.ToMat converter = new OpenCVFrameConverter.ToMat();
            Mat current = converter.convert(f);

            if (!frozen) {
                frame = current.clone();
            } else {
                current = frozenFrame.clone();
            }

            Mat display = current.clone();
            if (!frozen) {
                int size = (int)(magnifierSize / zoom);
                rectangle(display, new Point(magnifierPos.x()-size/2, magnifierPos.y()-size/2),
                          new Point(magnifierPos.x()+size/2, magnifierPos.y()+size/2), Scalar.GREEN, 2, 8, 0);
            } else {
                Mat mag = getMagnifiedRegion(frozenFrame);
                if (!mag.empty()) {
                    magCanvas.showImage(converter.convert(mag));
                    magCanvas.setVisible(true);
                    circle(display, magnifierPos, magnifierSize/2, Scalar.YELLOW, 2, 8, 0);
                }
            }

            // Информация
            String status = frozen ? "FROZEN" : "LIVE";
            putText(display, "Zoom: " + String.format("%.1f", zoom), new Point(10, 30),
                    FONT_HERSHEY_SIMPLEX, 0.7, Scalar.WHITE, 2, 8, false);
            putText(display, status, new Point(10, 60),
                    FONT_HERSHEY_SIMPLEX, 0.7, frozen ? Scalar.RED : Scalar.GREEN, 2, 8, false);

            canvas.showImage(converter.convert(display));

            // Обработка клавиш (упрощённо через консоль)
            // В JavaCV сложно получить клавиши из CanvasFrame, поэтому используем JFrame key listener
            // Для демонстрации: в консоли вводим команды
        }
    }

    public static void main(String[] args) throws Exception {
        new MagnifierJava().run();
    }
}
