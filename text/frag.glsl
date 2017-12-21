#version 410

uniform sampler2D tex;
in vec2 tex_pos;

void main() {
    gl_FragColor = vec4(1, 1, 1, texture2d(tex, tex_pos).r) * vec4(1, 1, 1, 1);
}
