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
	while(!glfwWindowShouldClose(app->ogl.window_hnd))
	{
		opencl_execute_kernel(app->particle_count, &app->ocl, &app->ogl, app->ocl.kupdate);
		glClearColor(0.0f, 0.0f, 0.0f, 1.0f);
		glClear(GL_COLOR_BUFFER_BIT);
		glBindVertexArray(app->ogl.vao);
		glDrawArrays(GL_POINTS, 0, app->particle_count);
		glfwSwapBuffers(app->ogl.window_hnd);
		glfwPollEvents();    
	}
}


int	main(void)
{
	t_app		app;

	app.particle_count = 3000000;
	init_glfw(&app);

	program_init(&app);
	buffer_init(&app);

	opencl_init(&app.ocl, &app.ogl);

	opencl_execute_kernel(app.particle_count, &app.ocl, &app.ogl, app.ocl.kinit);

	render(&app);
	return (1);
}