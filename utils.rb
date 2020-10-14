def create_shared_context(platform)
	if RUBY_PLATFORM =~ /cygwin|mswin|mingw|bccwin|wince|emx/
		options = {
			properties: [
				OpenCL::GL_CONTEXT_KHR, wglGetCurrentContext(),
				OpenCL::WGL_HDC_KHR, wglGetCurrentDC(),
				OpenCL::CONTEXT_PLATFORM, platform,
				0
			],
			user_data: nil
		}
	elsif RUBY_PLATFORM =~ /linux/
		options = {
			properties: [
				OpenCL::GL_CONTEXT_KHR, glXGetCurrentContext(),
				OpenCL::GLX_DISPLAY_KHR, glXGetCurrentDisplay(nil),
				OpenCL::CONTEXT_PLATFORM, platform,
				0
			],
			user_data: nil
		}
	elsif RUBY_PLATFORM =~ /darwin/
		ctx = CGLGetCurrentContext()
		share_group = CGLGetShareGroup(ctx)
		options = {
			properties: [
				CL_CONTEXT_PROPERTY_USE_CGL_SHAREGROUP_APPLE,
				share_group,
				0
			],
			user_data: nil
		}
	else
		exit()
	end
	return (options)
end