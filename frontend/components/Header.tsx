"use client"
import Link from 'next/link';
import Image from 'next/image';
import { useState } from 'react';

const NavBar = () => {
  const [isLoggedIn, setIsLoggedIn] = useState(true); // Replace this with your actual login check

  return (
    <nav className="flex justify-between items-center p-4 bg-gray-800 text-white">
      <div className={`flex items-center ${isLoggedIn ? 'justify-start' : 'justify-center'} flex-1`}>
        <Link href="/">
            logo here
        </Link>
      </div>

      {isLoggedIn && (
        <>
          <div className="flex-1 flex justify-center space-x-4">
            <Link href="/">
              Home
            </Link>
            <Link href="/about">
              <>About</>
            </Link>
            <Link href="/contact">
              <>Contact</>
            </Link>
          </div>
          <div className="flex-1 flex justify-end">
            <Image src="/profile.jpg" alt="profile" width={50} height={50} className="rounded-full" />
          </div>
        </>
      )}
    </nav>
  );
};

export default NavBar;
