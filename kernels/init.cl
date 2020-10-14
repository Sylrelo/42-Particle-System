__kernel void	initParticles(__global float3	*particles)
{
	int		gid = get_global_id(0);

	float	x, y, z, n, r, oa, teta;
	int		oa_index, teta_index, n_index;
	int		i = gid;

	particles[gid].xyz = (float3)(0.5, 0.5, 0.5);
}