/* ************************************************************************** */
/*                                                                            */
/*                                                        :::      ::::::::   */
/*   main.c                                             :+:      :+:    :+:   */
/*                                                    +:+ +:+         +:+     */
/*   By: slopez <slopez@student.42lyon.fr>          +#+  +:+       +#+        */
/*                                                +#+#+#+#+#+   +#+           */
/*   Created: 2020/10/16 18:01:42 by slopez            #+#    #+#             */
/*   Updated: 2020/10/16 18:01:42 by slopez           ###   ########lyon.fr   */
/*                                                                            */
/* ************************************************************************** */

#include <stdlib.h>
#include <stdio.h>
#include "global.h"
#include "GL/gl.h"
#include <math.h>

void		mat_action(mat4f out, const char action)
{
	if (action == 'i') {
		out[0][0] = 1.f;
		out[0][1] = 0.f;
		out[0][2] = 0.f;
		out[0][3] = 0.f;
		out[1][0] = 0.f;
		out[1][1] = 1.f;
		out[1][2] = 0.f;
		out[1][3] = 0.f;
		out[2][0] = 0.f;
		out[2][1] = 0.f;
		out[2][2] = 1.f;
		out[2][3] = 0.f;
		out[3][0] = 0.f;
		out[3][1] = 0.f;
		out[3][2] = 0.f;
		out[3][3] = 1.f;
	}
}

void		mat_persp(mat4f result, float fov, float aspect)
{
	float	tan_half_fov;
	float	near_plane;
	float	far_plane;

	near_plane = 0.1;
	far_plane = 1000.0;
	mat_action(result, 'i');
	tan_half_fov = tanf(fov / 2.0f);
	result[0][0] = 1.0f / (aspect * tan_half_fov);
	result[1][1] = 1.0f / (tan_half_fov);
	result[2][2] = (far_plane + near_plane) / (far_plane - near_plane) * -1;
	result[2][3] = -1.0f;
	result[3][2] = (2.0f * far_plane * near_plane) / (far_plane - near_plane) * -1;
}

void	die(const char *msg)
{
	printf("%s\n", msg);
	exit(1);
}

int		init_glfw(t_app *app)
{
	if (!glfwInit())
		die("Failed to initialize GLFW\n");
	glfwWindowHint(GLFW_CONTEXT_VERSION_MAJOR, 4);
	glfwWindowHint(GLFW_CONTEXT_VERSION_MINOR, 2);
	glfwWindowHint(GLFW_OPENGL_FORWARD_COMPAT, GL_TRUE);
	glfwWindowHint(GLFW_OPENGL_PROFILE, GLFW_OPENGL_CORE_PROFILE);
	app->ogl.window_hnd = glfwCreateWindow(1600, 900, "ParticleSystem", NULL, NULL);
	glfwMakeContextCurrent(app->ogl.window_hnd);
	if (gl3wInit())
		die ("Failed to initialize GL3W\n");
}

int		render(t_app *app)
{
	glUseProgram(app->ogl.program);

	mat4f	perspective;
	GLuint persp = glGetUniformLocation(app->ogl.program, "persp");

	//glEnable(GL_DEPTH_TEST);
	mat_persp(perspective, 60, 1600.0/900.0);
	glUniformMatrix4fv(persp, 1, GL_FALSE, perspective[0]);

	while(!glfwWindowShouldClose(app->ogl.window_hnd))
	{
		opencl_execute_kernel(app->particle_count, &app->ocl, &app->ogl, app->ocl.kupdate);
		glfwPollEvents();

		glClearColor(0.0f, 0.0f, 0.0f, 1.0f);
		glClear(GL_COLOR_BUFFER_BIT | GL_DEPTH_BUFFER_BIT);

		glBindVertexArray(app->ogl.vao);
		glDrawArrays(GL_POINTS, 0, app->particle_count);
		glfwSwapBuffers(app->ogl.window_hnd);
		glFlush();
	}
}


int	main(void)
{
	t_app		app;

	app.particle_count = 3000000;
	init_glfw(&app);
	opengl_init_program(&app);
	opengl_init_buffer(&app);
	opencl_init(&app.ocl, &app.ogl);
	opencl_execute_kernel(app.particle_count, &app.ocl, &app.ogl, app.ocl.kinit);
	clReleaseKernel(app.ocl.kinit);
	render(&app);
	return (1);
}