
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
	attr_accessor :kernel, :queue, :context, :kernel_init
	
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
			@platform 	= OpenCL.platforms.first
			@device 	= @platform.devices.first
			@context 	= OpenCL.create_context(@device, create_shared_context(@platform))
			@queue 		= OpenCL.create_command_queue(@context, @device)
			@prog 		= OpenCL.create_program_with_source(@context, init_kernel_src).build
		rescue => e
			puts "[OpenCL] Error Code : %s" % [e]
			begin puts(@prog.build_log) ; rescue ; end
			exit()
		end

		@kernel_init = OpenCL.create_kernel(@prog, "initParticles")

		begin
			@cl_buff = OpenCL.create_from_gl_buffer(@context, ogl.vbo_ptr)
		rescue => e
			puts("[OpenCL] Buffer creation from OpenGL failed")
			puts(e)
			exit()
		end

		OpenCL.enqueue_acquire_gl_objects(@queue, @cl_buff)
		OpenCL.set_kernel_arg(@kernel_init, 0, @cl_buff)

		OpenCL.enqueue_ndrange_kernel(@queue, @kernel_init, [main.count])
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
		glfwWindowHint(GLFW_CONTEXT_VERSION_MINOR, 2)
		glfwWindowHint(GLFW_OPENGL_FORWARD_COMPAT, GL_TRUE)
		glfwWindowHint(GLFW_OPENGL_PROFILE, GLFW_OPENGL_CORE_PROFILE)
		glfwWindowHint(GLFW_VISIBLE, GLFW_FALSE)

		@window = glfwCreateWindow(width, height, "Particles System", nil, nil)
		glfwMakeContextCurrent(window)
		ratio = width.to_f / height.to_f
		glViewport(0, 0, width, height)
		#glMatrixMode(GL_PROJECTION)
		#glLoadIdentity()
		#gluPerspective(30.0, ratio, 1.0, 1000.0)
		
		#glMatrixMode(GL_MODELVIEW)
	end
end

class OGL
	attr_accessor :vbo_ptr, :main

	def initialize(main)
		@main = main

		vbo = ' ' * 4
		vao = ' ' * 4

		glGenBuffers(1, vbo)
		glGenVertexArrays(1, vao)
		@vao_ptr = vao.unpack('L')[0]
		@vbo_ptr = vbo.unpack('L')[0]

		glBindVertexArray(@vao_ptr)
		glBindBuffer(GL_ARRAY_BUFFER, @vbo_ptr)
		glBufferData(GL_ARRAY_BUFFER, (Fiddle::SIZEOF_FLOAT * 4) * main.count, nil, GL_DYNAMIC_DRAW)
		glEnableVertexAttribArray(0)
		glVertexAttribPointer(0, 3, GL_FLOAT, GL_FALSE, 0, nil)
		self.init_programs()
	end

	def create_shader(src, type)
		shader = glCreateShader(type)
		glShaderSource(shader, 1, [src].pack('p'), [src.size].pack('I'))
		glCompileShader(shader)

		ret_value_buffer = ' ' * 4
		glGetShaderiv(shader, GL_COMPILE_STATUS, ret_value_buffer);
		ret_value = ret_value_buffer.unpack('L')[0]

		if ret_value == 0
			puts "Error in compiling shader"
			exit()
		end
		return (shader)
	end

	def init_programs()
		vertex_file = "kernels/vertex.glsl"
		fragment_file = "kernels/fragment.glsl"

		begin
			vertex_src = File.read(vertex_file)
		rescue
			puts "Can't open %s" % [vertex_file]
			exit()
		end
		begin
			fragment_src = File.read(fragment_file)
		rescue
			puts "Can't open %s" % [fragment_file]
			exit()
		end

		@program = glCreateProgram()

		vertex_shader 	= self.create_shader(vertex_src, GL_VERTEX_SHADER)
		fragment_shader = self.create_shader(fragment_src, GL_FRAGMENT_SHADER)

		glAttachShader(@program, vertex_shader)
		glAttachShader(@program, fragment_shader)
		glLinkProgram(@program)
		glDetachShader(@program, vertex_shader)
		glDetachShader(@program, fragment_shader)
		glDeleteShader(vertex_shader)
		glDeleteShader(fragment_shader)
		glUseProgram(@program)

	end

	def Display()
		glfwShowWindow(@main.window)

		while glfwWindowShouldClose( @main.window ) == 0
			glClear(GL_COLOR_BUFFER_BIT | GL_DEPTH_BUFFER_BIT)
			glBindVertexArray(@vao_ptr);
			glDrawArrays(GL_POINTS, 0, @main.count)
			glfwSwapBuffers( @main.window )
			glfwPollEvents()
		end
	end

end

if __FILE__ == $0
	main = Main.new()
	main.width = 1600
	main.height = 900
	main.count = 3000000
	main.initGlfw()
	ogl = OGL.new(main)
	ocl = OCL.new(main, ogl)
	ogl.Display()
end