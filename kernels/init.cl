__kernel void addition(  float2 alpha, __global const float *x, __global float *y) {
  size_t ig = get_global_id(0);
  y[ig] = (alpha.s0 + alpha.s1 + x[ig])*0.3333333333333333333f;
}