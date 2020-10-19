__kernel void	initParticles(__global float3 *particles, __global float3 *velocity)
{
	int				i = get_global_id(0);
	
	//if (i >= 2200000)
	//	printf("Current : %d %d\n", i, get_global_size(0));

	/*if (get_global_id(0) == 0) {
		printf("Local Sizes: %d %d \n", get_local_size(0));
		printf("Global Sizes: %d %d \n", get_global_size(0));
	}

	if (get_group_id(0) == 0){

		printf("Work Item Id: %d %d \n", get_local_id(0), get_local_id(1));
	}
	*/


	uint	subDivCount = cbrt((float)3000000);
	float2	delta = (float2)(M_PI * 2 / subDivCount, M_PI / subDivCount);
	float	radiusDelta = 5.f / subDivCount;

	uint	x = fmod(i, (float)subDivCount);
	uint	y = fmod((float)i / subDivCount, (float)subDivCount);
	uint	r = i / (subDivCount * subDivCount);

	float radius = radiusDelta * r;

	float2	offset = (r % 2 == 0) ? delta / 2 : (float2)(0);

	particles[i].x = radius * sin(delta.x * y + offset.x) * sin(delta.y * x + offset.y);
	particles[i].y = radius * cos(delta.x * y + offset.x);
	particles[i].z = (radius * sin(delta.x * y + offset.x) * cos(delta.y * x + offset.y));

	velocity[i].x = 0.0f;
	velocity[i].y = 0.0f;
	velocity[i].z = 0.0f;
}

float3			get_gp_effect(float3 particles)
{
	float3 op = {0.0, 0.0, 5};

	float3	direction = op.xyz - particles.xyz;
	float	distance = length(direction);
	float3	velocity = normalize(direction) * (1.f / (distance * 3));

	float scalarVelocity = length(velocity);

	if (scalarVelocity > 2.0f)
		velocity = normalize(velocity) * 2.0f;

	return velocity;
}

__kernel void	updateParticles(__global float3 *particles, __global float3 *velocity)
{
	int				i = get_global_id(0);

	if (i == 0)
		printf("Connard\n");

	if (length(velocity[i]) > 0.5f) {
		velocity[i] /= 1.04f;
	}

	velocity[i] 	+= (get_gp_effect(particles[i]));
	particles[i] 	+= velocity[i] / 1000.f;
}


