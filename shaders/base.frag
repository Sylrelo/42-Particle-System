#version 410

out vec4  outputColor;

in vec3   velocity;


vec3 speed_to_color(float speed, float min_speed, float max_speed) {
    float normalized_speed = clamp((speed - min_speed) / (max_speed - min_speed), 0.0, 1.0);

    vec3 orange = vec3(1.0, 0.5, 0.0); 
    vec3 blue = vec3(0.4, 0.4, 1.0);  

    return mix(orange, blue, normalized_speed);
}

void main()
{

    // if (fragColor.r == 0 && fragColor.g == 0 && fragColor.b == 0) {
    //     outputColor = vec4(1.0, 1.0, 1.0, 0.35);
    // }
    // else {
    //     outputColor = vec4(fragColor.rgb, 0.55);
    // }

    vec3 color = vec3(0.5, 0.5, 0.5);

    // color += velocity * 100.0;
    color = speed_to_color(length(velocity), 0.0, 0.02);

    color += velocity * 5.0;
    

    outputColor = vec4(color, 0.8);
}