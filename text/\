#version 410

in vec3 pos;
out vec2 tex_pos;

void main() {
    gl_Position = vec4(pos.xy, 0, 1);
    tex_pos = pos.zw;
}
