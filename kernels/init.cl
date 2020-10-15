__kernel void	initParticles(__global float3	*particles)
{
	int		gid = get_global_id(0);

	particles[gid].xyz = (float3)(0.5, 0.5, 0.5);
}