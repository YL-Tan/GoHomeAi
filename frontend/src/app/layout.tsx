import type { Metadata } from "next";
import { Inter } from "next/font/google";
import Link from "next/link";
import "./globals.css";
import { WebSocketProvider } from "@/context/WebSocketProvider";

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "GoHome AI Dashboard",
  description: "Monitor system metrics and job statuses in real time",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className={`${inter.className} bg-gray-900 text-white min-h-screen`}>
        {/* Navigation Bar */}
        <nav className="p-4 bg-gray-800 shadow-md flex justify-between items-center">
          <h1 className="text-xl font-bold">GoHome AI</h1>
          <div className="space-x-6">
            <Link href="/" className="hover:text-gray-400">Home</Link>
            <Link href="/dashboard" className="hover:text-gray-400">Dashboard</Link>
            <Link href="/metrics" className="hover:text-gray-400">Metrics</Link>
          </div>
        </nav>
        {/* Main Content */}
        <WebSocketProvider> <main className="p-6 flex-1">{children}</main></WebSocketProvider>
      </body>
    </html>
  );
}
