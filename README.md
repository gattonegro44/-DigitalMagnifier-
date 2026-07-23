🔍 DigitalMagnifier — цифровая лупа с заморозкой кадра
Интерактивный инструмент для увеличения и детального рассмотрения изображений с веб-камеры.
Позволяет заморозить кадр, перемещать увеличительное стекло, регулировать увеличение и сохранять результат.
Реализован на 7 языках программирования для демонстрации подходов к обработке видео и интерактивным интерфейсам.

https://img.shields.io/github/repo-size/yourname/digitalmagnifier
https://img.shields.io/github/stars/yourname/digitalmagnifier?style=social
https://img.shields.io/badge/License-MIT-blue.svg

🧠 Концепция
DigitalMagnifier — это приложение, которое превращает веб-камеру в цифровую лупу. Оно позволяет:

✅ Просматривать видео с веб-камеры в реальном времени.

✅ Замораживать кадр — останавливать видео для детального изучения.

✅ Перемещать лупу по замороженному кадру с помощью мыши или клавиш.

✅ Регулировать увеличение (зум) от 1x до 10x.

✅ Отображать увеличенную область в отдельном окне или на том же экране.

✅ Сохранять замороженный кадр или увеличенную область в файл (PNG).

✅ Управлять через интуитивный интерфейс (кнопки, слайдеры, горячие клавиши).

🚀 Как запустить
Каждая версия использует соответствующие библиотеки. Инструкции по установке и запуску:

Python
bash
pip install opencv-python numpy
python magnifier_python.py
C++
bash
# Требуется OpenCV (sudo apt install libopencv-dev)
g++ -std=c++17 magnifier_cpp.cpp -o magnifier `pkg-config --cflags --libs opencv4`
./magnifier
Java
bash
# Требуется JavaCV (скачать jar)
javac -cp .:javacv.jar MagnifierJava.java
java -cp .:javacv.jar MagnifierJava
C# (.NET Core)
bash
dotnet add package OpenCvSharp4.Windows
dotnet run
Go
bash
# Требуется OpenCV для Go (go-opencv)
go get -u github.com/hybridgroup/go-opencv/...
go run magnifier_go.go
Rust
bash
# Требуется OpenCV для Rust (opencv-rust)
cargo build --release
./target/release/magnifier_rs
JavaScript (браузер)
bash
# Откройте magnifier_js.html в современном браузере
🧩 Пример использования
text
🔍 Цифровая лупа (заморозка кадра)
[Пробел] — заморозить/возобновить видео
[Мышь] — перемещать лупу
[Колесо] — изменение увеличения
[S] — сохранить кадр
[Q] — выход
📦 Содержимое репозитория
Файл	Язык	Особенности
magnifier_python.py	Python	OpenCV, обработка мыши, слайдер зума, сохранение PNG
magnifier_cpp.cpp	C++	OpenCV, окно с лупой, горячие клавиши
MagnifierJava.java	Java	JavaCV, Swing-интерфейс, кнопки управления
MagnifierCSharp.cs	C#	OpenCvSharp, WPF-интерфейс, слайдер
magnifier_go.go	Go	GoCV, консольный интерфейс с выбором области
magnifier_rs.rs	Rust	opencv-rust, консольное управление
magnifier_js.html	JavaScript	WebRTC + Canvas, браузерный интерфейс с лупой
🔮 Расширенные функции
Поддержка нескольких камер (переключение).

Фильтры (оттенки серого, инверсия) для увеличенной области.

Автоматическое сохранение при каждом замораживании.

Режим «картинка в картинке» (увеличенная область поверх основного видео).

📜 Лицензия
MIT — свободно используйте, модифицируйте и распространяйте.

🤝 Вклад
Приветствуются пул-реквесты с улучшениями, поддержкой новых платформ и расширением функциональности.

⭐ Если проект помогает вам разглядеть детали — поставьте звёздочку!

