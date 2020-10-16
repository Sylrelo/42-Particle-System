#version 400 core
layout (location = 0) in vec3 aPos;
  
out vec4 vertexColor;

void main()
{
    float testx = aPos.x;
    float testy = aPos.y;
    float testz = aPos.z;

    if (testz < 0.0)
        testz = 0.0f;
    if (testz > 1.0)
        testz = 1.0f;

    if (testx < 0.0)
        testx = 0.0f;
    if (testx > 1.0)
        testx = 1.0f;

    if (testy < 0.0)
        testy = 0.0f;
    if (testy > 1.0)
        testy = 1.0f;

    gl_Position = vec4(aPos, 1.0);
    vertexColor = vec4(testx, testy, testz, 1.0);
}