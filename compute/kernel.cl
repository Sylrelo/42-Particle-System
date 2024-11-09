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

// TODO
__kernel void gravitateParticles(__global float4 *particles, __global float4 *velocity) {
  int i = get_global_id(0);

  float3 TMP_GRAVITATION_POINT = (float3) {0, 0.52, -2.5};
  float GRAVITY_CONST = 0.0001f;
  float DELTA_TIME = 2.1;

  float3 direction = TMP_GRAVITATION_POINT - particles[i].xyz;
  float distance = length(direction) + 1e-6f;
  direction /= distance;
  float force = GRAVITY_CONST / (distance * distance);


  velocity[i].xyz += direction * force * DELTA_TIME;
  particles[i].xyz += velocity[i].xyz * DELTA_TIME;
}

__kernel void initParticles(__global float4 *particles, __global float4 *velocity) {
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

  //TMP
  velocity[i] = (float4) {0.0, 0.001, 0.002, 0 };
  // colors[i] = (float4) { 1, 1, 1, 0 };
}