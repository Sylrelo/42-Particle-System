gcc -L./libs -L"C:\Program Files\NVIDIA GPU Computing Toolkit\CUDA\v11.1\lib\Win32" src/*.c libs/gl3w.c -I./includes -lOpenCL -lglfw3dll -lopengl32 -o build/particlesystem.exe
build\particlesystem.exe