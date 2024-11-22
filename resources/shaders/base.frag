#version 410

out vec4  outputColor;

in vec3     velocity;
in vec3     position;

void main()
{
    float speed = length(velocity);
    vec3 color = vec3(0.5, 0.5, 0.5);

    // TESTS
    vec3 fastColor = vec3(0.8, 0.8, 1.0);
    vec3 slowColor = vec3(0.5, 0.3, 1.0);

    float maxSpeed = 0.05;
    float normalizedSpeed = clamp(speed / maxSpeed, 0.0, 1.0);


    color = mix(slowColor, fastColor, normalizedSpeed);
    color += vec3(
        0.3 * sin(position.x * 0.1),
        0.3 * sin(position.y * 0.1 + 1.0),
        0.3 * sin(position.z * 0.1 + 2.0)
    );

    color = clamp(color, 0.0, 1.0);

    // color += velocity * 100.0;
    // color = speed_to_color(length(velocity), 0.0, 0.02);
    // color += velocity * 2.0;
    

    outputColor = vec4(color, 0.5);
}