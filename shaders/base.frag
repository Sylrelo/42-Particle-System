#version 410

out vec4  outputColor;
// in vec3   fragColor;

void main()
{

    // if (fragColor.r == 0 && fragColor.g == 0 && fragColor.b == 0) {
    //     outputColor = vec4(1.0, 1.0, 1.0, 0.35);
    // }
    // else {
    //     outputColor = vec4(fragColor.rgb, 0.55);
    // }

    outputColor = vec4(1.0);
}