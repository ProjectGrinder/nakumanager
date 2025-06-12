import type { Metadata } from "next";
import { Noto_Sans } from "next/font/google"
import "./globals.css";

const notoSans = Noto_Sans({
  weight: ["700", '500', "400", "300"],
  subsets: ["latin", "latin-ext"],
  variable: '--font-noto-sans',
  display: "swap"
})

export const metadata: Metadata = {
  title: "NakuManager",
  description: "Open source manager application",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body
        className={`${notoSans.className} antialiased`}
      >
        {children}
      </body>
    </html>
  );
}
