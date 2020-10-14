
require 'opencl_ruby_ffi'
require 'narray_ffi'
require 'opengl'
require 'glfw'
require 'glu'

require './utils.rb'

include GLFW
include GLU
include OpenGL

OpenGL.load_lib()
GLFW.load_lib()
GLU.load_lib()

class OCL
	attr_accessor :kernel, :queue, :context

	def initialize(main, ogl)
		@main = main
		@ogl = ogl
		init_kernel_file = "kernels/init.cl"

		begin
			init_kernel_src = File.read(init_kernel_file)
		rescue
			puts "Can't open %s" % [init_kernel_file]
			exit()
		end

		begin
			@platform = OpenCL::platforms[0]
		rescue => e
			puts("[OpenCL] Platform selection failed")
			puts(e)
			exit()
		end

		begin
			@device = @platform.devices[0]
		rescue => e
			puts("[OpenCL] Device selection failed")
			puts(e)
			exit()
		end

		begin
			@context = OpenCL.create_context(@device, create_shared_context(@platform))
		rescue => e
			puts("[OpenCL] Context creation failed")
			puts(e)
			exit()
		end

		begin
			@queue = @context.create_command_queue(@device, :properties => OpenCL::CommandQueue::PROFILING_ENABLE)
		rescue => e
			puts("[OpenCL] Command queue creation failed")
			puts(e)
			exit()
		end

		begin
			@prog = @context.create_program_with_source( init_kernel_src )
		rescue => e
			puts("[OpenCL] Program creation failed")
			puts(e)
			exit()
		end

		begin
			@kernel = @prog.build
		rescue => e
			puts("[OpenCL] Build failed")
			puts(e)
			puts(@prog.build_log)
			exit()
		end

		begin
			@cl_buff = OpenCL.create_from_gl_buffer(@context, ogl.mesh_vbo[0])
		rescue => e
			puts("[OpenCL] Buffer creation from OpenGL failed")
			puts(e)
			exit()
		end

		begin
			OpenCL.enqueue_acquire_gl_objects(@queue, @cl_buff)
		rescue => e
			puts("[OpenCL] Acquire GL Buffer failed")
			puts(e)
			exit()
		end
		OpenCL.enqueue_release_gl_objects(@queue, @cl_buff)

		OpenCL.finish(@queue)
	end


end

class Main
	attr_reader :window
	attr_accessor :width, :height, :count

	def initialize()
		puts("Starting Particle System...")
	end

	def initGlfw()
		glfwInit()
		glfwWindowHint(GLFW_RESIZABLE, GLFW_FALSE)
		glfwWindowHint(GLFW_SAMPLES, 4)
		glfwWindowHint(GLFW_CONTEXT_VERSION_MAJOR, 4)
		glfwWindowHint(GLFW_CONTEXT_VERSION_MINOR, 0)
		glfwWindowHint(GLFW_OPENGL_FORWARD_COMPAT, GL_TRUE)
		glfwWindowHint(GLFW_OPENGL_PROFILE, GLFW_OPENGL_CORE_PROFILE)
		glfwWindowHint(GLFW_VISIBLE, GLFW_FALSE)

		@window = glfwCreateWindow( width, height, "Particles System", nil, nil )
		glfwMakeContextCurrent( window )
		ratio = width.to_f / height.to_f
		glViewport(0, 0, 1280, 720)
		glMatrixMode(GL_PROJECTION)
		glLoadIdentity()
		gluPerspective(45.0, ratio, 1.0, 1000.0)
		glMatrixMode(GL_MODELVIEW)
	end

	def Display()
		glfwShowWindow(@window)
		while glfwWindowShouldClose( window ) == 0
			glClear(GL_COLOR_BUFFER_BIT | GL_DEPTH_BUFFER_BIT)
			#glLoadIdentity()
			#glTranslatef(0.0, 0.0, -50.0)
			#glRotatef(glfwGetTime() * 50.0, 0.0, 1.0, 0.0)
			#teapot.render
			glfwSwapBuffers( window )
			glfwPollEvents()
		end
	end
end

class OGL
	attr_accessor :mesh_vbo

	def initialize(main)
		@main = main

		vbo = ' ' * 4
		glGenBuffers(1, vbo)
		@mesh_vbo = vbo.unpack('L')

		glBindBuffer(GL_ARRAY_BUFFER, @mesh_vbo[0])
		glBufferData(GL_ARRAY_BUFFER, Fiddle::SIZEOF_FLOAT * main.count, nil, GL_STATIC_DRAW)
		#glBindBuffer(GL_ARRAY_BUFFER, 0)
		glEnableVertexAttribArray(0)
		glVertexAttribPointer(0, 1, GL_FLOAT, GL_FALSE, 0, nil)
		
	end
end

if __FILE__ == $0
	main = Main.new()
	main.width = 1280
	main.height = 720
	main.count = 3000000

  	main.initGlfw()

	ogl = OGL.new(main)
	ocl = OCL.new(main, ogl)

	main.Display()
end