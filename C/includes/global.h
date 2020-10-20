/* ************************************************************************** */
/*                                                                            */
/*                                                        :::      ::::::::   */
/*   global.h                                           :+:      :+:    :+:   */
/*                                                    +:+ +:+         +:+     */
/*   By: slopez <slopez@student.42lyon.fr>          +#+  +:+       +#+        */
/*                                                +#+#+#+#+#+   +#+           */
/*   Created: 2020/10/16 18:02:13 by slopez            #+#    #+#             */
/*   Updated: 2020/10/16 18:02:13 by slopez           ###   ########lyon.fr   */
/*                                                                            */
/* ************************************************************************** */

#define CL_TARGET_OPENCL_VERSION 120
#include "GL/gl3w.h"
#include "GLFW/glfw3.h"
#include "CL/opencl.h"

typedef GLfloat	mat4f[4][4];
typedef struct			s_ogl
{
	GLFWwindow			*window_hnd;
	unsigned int		program;
	unsigned int		vbo_pos;
	unsigned int		vao;
}						t_ogl;

typedef struct			s_ocl
{
	cl_context			ctx;
	cl_device_id		device;
	cl_kernel			kinit;
	cl_kernel			kupdate;
	cl_mem				buff_pos;
	cl_mem				buff_velo;
	cl_command_queue	queue;
}						t_ocl;

typedef struct			s_app
{
	t_ogl				ogl;
	t_ocl				ocl;
	int					particle_count;

}						t_app;

void	die(const char *msg);


void	opengl_init_buffer(t_app *app);
void	opengl_init_program(t_app *app);

void	opencl_init(t_ocl *ocl, t_ogl *ogl);
void	opencl_execute_kernel(int particle_count, t_ocl *ocl, t_ogl *ogl, cl_kernel kernel);
void	opencl_error(int err);