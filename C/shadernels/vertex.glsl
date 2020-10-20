#version 400 core

layout (location = 0) in vec3 aPos;

uniform mat4	persp;

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

    gl_Position = persp * vec4(aPos, 1.0);
    vertexColor = vec4(0.5, 0.5, testz, 1.0);
}