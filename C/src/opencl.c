/* ************************************************************************** */
/*                                                                            */
/*                                                        :::      ::::::::   */
/*   opencl.c                                           :+:      :+:    :+:   */
/*                                                    +:+ +:+         +:+     */
/*   By: slopez <slopez@student.42lyon.fr>          +#+  +:+       +#+        */
/*                                                +#+#+#+#+#+   +#+           */
/*   Created: 2020/10/19 22:48:51 by slopez            #+#    #+#             */
/*   Updated: 2020/10/19 22:48:51 by slopez           ###   ########lyon.fr   */
/*                                                                            */
/* ************************************************************************** */

#include "global.h"
#include <stdlib.h>
#include <stdio.h>

cl_context_properties	*create_shared_context(cl_platform_id platform)
{
		static cl_context_properties properties[7];

		properties[0] = CL_GL_CONTEXT_KHR;
		properties[1] = (cl_context_properties)wglGetCurrentContext();
		properties[2] = CL_WGL_HDC_KHR;
		properties[3] = (cl_context_properties)wglGetCurrentDC();
		properties[4] = CL_CONTEXT_PLATFORM;
		properties[5] = (cl_context_properties)platform;
		properties[6] = 0;
		return (properties);
}

void	opencl_error(int err)
{
	if (err)
	{
		printf("An OpenCL error %d occured\n", err);
		exit(1);
	}
}

void	opencl_execute_kernel(int particle_count, t_ocl *ocl, t_ogl *ogl, cl_kernel kernel)
{
	cl_mem		buffers[2];
	int			err;
	size_t		local_work_size;

	local_work_size = 200;
	buffers[0] = ocl->buff_pos;
	buffers[1] = ocl->buff_velo;
	glFinish();
	err = clSetKernelArg(kernel, 0, sizeof(GL_FLOAT_VEC3), (void *)&ocl->buff_pos);
	opencl_error(err);
	err = clSetKernelArg(kernel, 1, sizeof(GL_FLOAT_VEC3), (void *)&ocl->buff_velo);
	opencl_error(err);
	err = clEnqueueAcquireGLObjects(ocl->queue, 2, buffers, 0, 0, 0);
	opencl_error(err);
	err = clEnqueueNDRangeKernel(ocl->queue, kernel, 1, 0, &particle_count, &local_work_size, 0, 0, 0);
	opencl_error(err);
	err = clEnqueueReleaseGLObjects(ocl->queue, 2, buffers, 0, 0, 0);
	opencl_error(err);
	clFinish(ocl->queue);
}

void	opencl_create_program(t_ocl *ocl)
{
	cl_program		program;
	FILE			*fp;
	char			*buffer;
	size_t			size;
	int				err;

	fp = fopen("shadernels/init.cl", "r");
	fseek(fp, 0, SEEK_END);
	size = ftell(fp);
	rewind(fp);
	buffer = malloc(size);
	if (!buffer)
		die("Can't allocate memory for shaders");
	memset(buffer, 0, size);
	fread(buffer, 1, size, fp);
	fclose(fp);
	program = clCreateProgramWithSource(ocl->ctx, 1, (const char **)&buffer,
		&size, &err);
	opencl_error(err);
	err = clBuildProgram(program, 1, &ocl->device, NULL, NULL, NULL);
	opencl_error(err);
	ocl->kinit = clCreateKernel(program, "initParticles", &err);
	opencl_error(err);
	ocl->kupdate = clCreateKernel(program, "updateParticles", &err);
	opencl_error(err);
	free(buffer);
}

void	opencl_init(t_ocl *ocl, t_ogl *ogl)
{
	cl_device_id		device_id;
	cl_platform_id		platform;
	int					err;

	err = clGetPlatformIDs(1, &platform, NULL);
	opencl_error(err);
	err = clGetDeviceIDs(platform, CL_DEVICE_TYPE_GPU, 1, &ocl->device, NULL);
	opencl_error(err);
	ocl->ctx = clCreateContext(create_shared_context(platform), 1, &ocl->device, NULL, NULL, &err);
	opencl_error(err);
	ocl->queue = clCreateCommandQueue(ocl->ctx, ocl->device, 0, &err);
	opencl_error(err);
	opencl_create_program(ocl);
	ocl->buff_pos = clCreateFromGLBuffer(ocl->ctx, CL_MEM_WRITE_ONLY, 
		ogl->vbo_pos, &err);
	opencl_error(err);
	ocl->buff_velo = clCreateFromGLBuffer(ocl->ctx, CL_MEM_WRITE_ONLY, 
		ogl->vbo_velo, &err);
	opencl_error(err);
}