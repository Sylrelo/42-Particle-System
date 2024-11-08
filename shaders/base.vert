#version 410 core

layout(location = 0) in vec4 aPosition;
layout(location = 1) in vec4 aVelocity;

uniform mat4 perspective;
// uniform mat4 model;
uniform mat4 camera;

// out vec3 fragColor;

out vec3 velocity;

void main(void)
{
    gl_Position = perspective * camera * vec4(aPosition.xyz, 1.0);
    velocity = aVelocity.xyz;
    // fragColor = vec4(1);
    // fragColor = aColors.rgb;
}