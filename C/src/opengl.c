/* ************************************************************************** */
/*                                                                            */
/*                                                        :::      ::::::::   */
/*   opengl.c                                           :+:      :+:    :+:   */
/*                                                    +:+ +:+         +:+     */
/*   By: slopez <slopez@student.42lyon.fr>          +#+  +:+       +#+        */
/*                                                +#+#+#+#+#+   +#+           */
/*   Created: 2020/10/17 19:08:22 by slopez            #+#    #+#             */
/*   Updated: 2020/10/17 19:08:22 by slopez           ###   ########lyon.fr   */
/*                                                                            */
/* ************************************************************************** */

#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include "global.h"
#include "GL/gl.h"

static void			shader_compile(unsigned int shader, char *buffer)
{
	int		success;
	char	log[1024];

	glShaderSource(shader, 1, (const char *const *)&buffer, NULL);
	glCompileShader(shader);
	glGetShaderiv(shader, GL_COMPILE_STATUS, &success);
	if (!success)
	{
		glGetShaderInfoLog(shader, 1024, NULL, log);
		printf("%s\n", log);
		exit(1);
	}
}

static unsigned int	shader_init(t_app *app, const char *file, unsigned int type)
{
	FILE			*fp;
	char			*buffer;
	long long int	size;
	unsigned int	shader;

	fp = fopen(file, "r");
	fseek(fp, 0, SEEK_END);
	size = ftell(fp);
	rewind(fp);
	buffer = malloc(size);
	if (!buffer)
		die("Can't allocate memory for shaders");
	memset(buffer, 0, size);
	fread(buffer, 1, size, fp);
	fclose(fp);
	shader = glCreateShader(type);
	shader_compile(shader, buffer);
	free(buffer);
	return (shader);
}

void				opengl_init_program(t_app *app)
{
	unsigned int	fragment;
	unsigned int	vertex;
	int				success;
	char			log[1024];

	vertex = shader_init(app, "shadernels/vertex.glsl", GL_VERTEX_SHADER);
	fragment = shader_init(app, "shadernels/fragment.glsl", GL_FRAGMENT_SHADER);
	app->ogl.program = glCreateProgram();
	glAttachShader(app->ogl.program, vertex);
	glAttachShader(app->ogl.program, fragment);
	glLinkProgram(app->ogl.program);
	glGetProgramiv(app->ogl.program, GL_LINK_STATUS, &success);
	if (!success)
	{
		glGetProgramInfoLog(app->ogl.program, 1024, NULL, log);
		printf("%s\n", log);
		exit(1);
	}
	glDeleteShader(vertex);
	glDeleteShader(fragment);
}

void				opengl_init_buffer(t_app *app)
{
	int				err;
	long unsigned	size;

	size = (sizeof(float) * 4) * app->particle_count;
	glGenVertexArrays(1, &app->ogl.vao);
	glGenBuffers(1, &app->ogl.vbo_pos);
	glBindVertexArray(app->ogl.vao);
	glBindBuffer(GL_ARRAY_BUFFER, app->ogl.vbo_pos);
	glBufferData(GL_ARRAY_BUFFER, size, NULL, GL_DYNAMIC_DRAW);
	glEnableVertexAttribArray(0);
	glVertexAttribPointer(0, 4, GL_FLOAT, GL_FALSE, 0, NULL);
	glFinish();
	if (glGetError())
	{
		printf("%d\n", err);
		exit(1);
	}
}
