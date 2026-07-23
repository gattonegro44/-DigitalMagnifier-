// MagnifierCSharp.cs — цифровая лупа с заморозкой кадра на C# (OpenCvSharp)

using System;
using OpenCvSharp;

class DigitalMagnifier
{
    private VideoCapture capture;
    private Mat frame, frozenFrame;
    private bool frozen = false;
    private Point magnifierPos = new Point(200, 200);
    private double zoom = 2.0;
    private int magnifierSize = 150;
    private string windowName = "Digital Magnifier";

    public DigitalMagnifier()
    {
        capture = new VideoCapture(0);
        if (!capture.IsOpened())
        {
            Console.WriteLine("Не удалось открыть камеру");
            Environment.Exit(1);
        }
        Cv2.NamedWindow(windowName);
        Cv2.SetMouseCallback(windowName, MouseCallback);
        Console.WriteLine("🔍 Цифровая лупа (заморозка кадра)");
        Console.WriteLine("Пробел — заморозить/возобновить");
        Console.WriteLine("Колесо мыши — изменение зума");
        Console.WriteLine("S — сохранить кадр");
        Console.WriteLine("Q — выход");
    }

    private void MouseCallback(MouseEventTypes evt, int x, int y, MouseEventFlags flags, IntPtr userdata)
    {
        if (evt == MouseEventTypes.Move)
            magnifierPos = new Point(x, y);
        else if (evt == MouseEventTypes.MouseWheel)
        {
            int delta = (int)flags; // приблизительно
            zoom += delta / 120.0 * 0.5;
            zoom = Math.Max(1.0, Math.Min(10.0, zoom));
        }
    }

    private Mat GetMagnifiedRegion(Mat src)
    {
        int x = magnifierPos.X, y = magnifierPos.Y;
        int regionSize = (int)(magnifierSize / zoom);
        int x1 = Math.Max(0, Math.Min(src.Cols, x - regionSize / 2));
        int y1 = Math.Max(0, Math.Min(src.Rows, y - regionSize / 2));
        int x2 = Math.Max(0, Math.Min(src.Cols, x + regionSize / 2));
        int y2 = Math.Max(0, Math.Min(src.Rows, y + regionSize / 2));
        if (x2 <= x1 || y2 <= y1) return null;
        Mat region = new Mat(src, new Rect(x1, y1, x2 - x1, y2 - y1));
        Mat magnified = new Mat();
        Cv2.Resize(region, magnified, new Size(magnifierSize, magnifierSize));
        Cv2.Circle(magnified, new Point(magnifierSize / 2, magnifierSize / 2), 5, Scalar.Green, 1);
        return magnified;
    }

    public void Run()
    {
        while (true)
        {
            if (!frozen)
            {
                frame = new Mat();
                capture.Read(frame);
                if (frame.Empty()) break;
            }
            else
            {
                frame = frozenFrame.Clone();
            }

            Mat display = frame.Clone();
            if (!frozen)
            {
                int size = (int)(magnifierSize / zoom);
                Cv2.Rectangle(display, new Point(magnifierPos.X - size / 2, magnifierPos.Y - size / 2),
                              new Point(magnifierPos.X + size / 2, magnifierPos.Y + size / 2), Scalar.Green, 2);
            }
            else
            {
                Mat mag = GetMagnifiedRegion(frozenFrame);
                if (mag != null)
                {
                    Cv2.ImShow("Magnifier View", mag);
                    Cv2.Circle(display, magnifierPos, magnifierSize / 2, Scalar.Yellow, 2);
                }
            }

            Cv2.PutText(display, $"Zoom: {zoom:F1}", new Point(10, 30), HersheyFonts.HersheySimplex, 0.7, Scalar.White, 2);
            string status = frozen ? "FROZEN" : "LIVE";
            Cv2.PutText(display, status, new Point(10, 60), HersheyFonts.HersheySimplex, 0.7,
                        frozen ? Scalar.Red : Scalar.Green, 2);

            Cv2.ImShow(windowName, display);

            char key = (char)Cv2.WaitKey(1);
            if (key == 'q') break;
            else if (key == ' ')
            {
                if (frozen)
                {
                    frozen = false;
                    Cv2.DestroyWindow("Magnifier View");
                }
                else if (frame != null)
                {
                    frozenFrame = frame.Clone();
                    frozen = true;
                }
            }
            else if (key == 's')
            {
                string filename = frozen ? "frozen_frame.png" : "live_frame.png";
                Cv2.ImWrite(filename, frozen ? frozenFrame : frame);
                Console.WriteLine($"Сохранено: {filename}");
            }
        }
        capture.Release();
        Cv2.DestroyAllWindows();
    }

    public static void Main()
    {
        new DigitalMagnifier().Run();
    }
}
