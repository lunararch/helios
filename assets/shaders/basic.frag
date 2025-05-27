#version 410 core

in vec2 TexCoord;
out vec4 FragColor;

uniform sampler2D texture1;
uniform vec4 color;

void main()
{
    FragColor = texture(texture1, TexCoord) * color;
}