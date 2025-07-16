type SidebarButtonProps = {
  children: React.ReactNode;
  onClick?: () => void;
};

export default function SidebarButton(props: SidebarButtonProps) {
  return (
    <button
      className="w-full text-sm text-left text-gray-400 p-2 rounded hover:bg-gray-700 cursor-pointer transition duration-200 truncate"
      onClick={props.onClick}
    >
      {props.children}
    </button>
  );
}
