#version 410 core

layout(location = 0) in vec4 aPosition;
layout(location = 1) in vec2 aTexCoord;

uniform mat4 perspective;
uniform mat4 camera;
uniform mat4 model;

out vec2 texCoords; 
out vec3 position;

void main(void)
{
  position = (model * vec4(aPosition.xyz, 1.0)).xyz;

  gl_Position = perspective * camera * model * vec4(aPosition.xyz, 1.0);
  texCoords = aTexCoord; 
}