export default function Home() {
  return (
    <main className="min-h-screen bg-gray-100 text-gray-900">
      <div className="grid grid-cols-1 lg:grid-cols-2 h-full border-primary border-t border-b">

        {/* Left Side */}
        <div className="flex flex-col justify-center p-5 sm:p-5 md:p-5 lg:p-16 border-r border-primary">
          {/* Company Name */}
          <div className="flex items-center justify-center w-full">
            <h1 className="text-5xl sm:text-6xl md:text-7xl lg:text-7xl xl:text-10xl font-bold leading-none text-center">
              YourCompany
            </h1>
          </div>
        </div>

        {/* Right Side */}
        <div className="sm:mt-0 md:mt-1 md:ml-5 md:mr-5 flex flex-col justify-start lg:pt-8">
          
          {/* Navbar */}
          <nav className="sm:mt-0 md:mt-0 text-2xl md:text:2xl lg:text-2xl grid grid-cols-2 border-b border-primary lg:pb-4 text-center font-medium">
            <div className="flex justify-center">
              <a href="#register" className="hover:underline">Register</a>
            </div>
            <div className="flex justify-center">
              <a href="#login" className="hover:underline">Login</a>
            </div>
          </nav>

          {/* Text */}
          <div className="mt-6 mb-6 ml-0 mr-0 border-primary pb-0">
            <h2 className="text-2xl font-light px-8 text-center">
              The all-in-one app for cooking, nutrition, and health. <br/>
              From discovering new recipes to managing your daily intake, allergens, and fitness goals. <br/>
              Everything you need in one place.<br/>
              Cook smarter, eat better, and take full control of your health journey.
            </h2>
          </div>

        </div>

      </div>
    </main>
  )
}