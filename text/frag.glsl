#version 410

uniform sampler2D tex;
in vec2 tex_pos;
out vec4 color;

void main() {
    color = vec4(1, 1, 1, texture2d(tex, tex_pos).r) * vec4(1, 1, 1, 1);
}
