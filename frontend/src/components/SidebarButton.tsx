export default function SidebarButton({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <button className="w-full text-sm text-left text-gray-400 p-2 rounded-sm hover:bg-gray-700 cursor-pointer transition duration-200">
      {children}
    </button>
  );
}
