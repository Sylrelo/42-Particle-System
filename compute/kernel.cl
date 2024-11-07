static ulong randu(ulong *seed) {
  ulong value = *seed * 1103515245 + 12345;

  *seed = value;
  return value;
}

static float rand_float(ulong *seed) { return (float)randu(seed) / ULONG_MAX; }

static float rand_float_in_range(ulong *seed, float a, float b) {
  float x = rand_float(seed);

  return (b - a) * x + a;
}

static float3 get_position_in_cube(ulong *seed) {
  float3 position;

  position.x = rand_float_in_range(seed, -1.0f, 1.0f);
  position.y = rand_float_in_range(seed, -1.0f, 1.0f);
  position.z = rand_float_in_range(seed, -1.0f, 1.0f);

  return position;
}

__kernel void initParticles(__global float4 *particles) {
  int i = get_global_id(0);

  ulong x = i / 1280.0;
  ulong y = i % 720;

  ulong seed = i * x * y;
  seed = randu(&seed);

  float3 rnd_position = get_position_in_cube(&seed);

  particles[i].x = rnd_position.x;
  particles[i].y = rnd_position.y;
  particles[i].z = rnd_position.z;
  // particles[i].x = float(x);
  // particles[i].y = float(y);
  // particles[i].z = 0;
  particles[i].w = 0;

  // printf("%f %f %f\n", particles[i].x, particles[i].y, particles[i].z);
  // particles[i].w = 0;

  // velocity[i] = (float4) { 0, 0, 0, 0 };
  // colors[i] = (float4) { 1, 1, 1, 0 };
}