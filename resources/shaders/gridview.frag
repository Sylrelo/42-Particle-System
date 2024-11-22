#version 410

out vec4  outputColor;

in vec2 texCoords;
in vec3 position;

uniform sampler2D gridTexture;

void main()
{
  float RATIO = 2000;

  vec4 texColor = texture(gridTexture, texCoords * RATIO);

  float dst = length(position);
  float fadeRatio = 1.0 - clamp(dst / 200, 0.05, 1.0);


  if (texColor.r < 0.9)
    discard;
  

  outputColor = vec4(vec3(0.6, 0.6, 0.8) * fadeRatio, 1.0);
  // outputColor = vec4(1.0);
  // outputColor = vec4(texture(gridTexture, texCoords * RATIO).xyz, 1.0);
}