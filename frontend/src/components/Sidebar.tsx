"use client";

import { useRouter } from "next/navigation";
import Link from "next/link";
import SidebarButton from "./SidebarButton";

export default function Sidebar() {
  const currentWorkspace = "Workspace 1";
  const teams = ["Team 1", "Team 2", "Team 3"];
  const router = useRouter();
  const handleLogout = () => {
    router.push("/login");
  };
  const createWorkspace = () => {
    router.push("/workspace-edit");
  };
  return (
    <div className="flex-column items-center justify-between p-2 border-r-1 border-gray-500 bg-gray-800 text-gray-200 w-60 h-screen relative">
      <div className="flex items-center justify-between py-2 mb-2">
        <span className="text-base font-semibold max-w-45 truncate">
          Username
        </span>
        <i
          className="fa-solid fa-right-from-bracket text-base"
          onClick={handleLogout}
        ></i>
      </div>
      <hr className="border-gray-500 mb-2"></hr>
      <div>
        <div className="flex flex-row items-center justify-between p-1 mb-2 text-gray-400 text-sm">
          <span className="max-w-45 truncate">{currentWorkspace}</span>
          <span className="w-6 h-6 text-center rounded-xl cursor-pointer hover:bg-gray-600 transition duration-200">
            ...
          </span>
        </div>
        <div className="flex flex-col ml-4">
          <SidebarButton>Manage members</SidebarButton>
          <SidebarButton>My issues</SidebarButton>
          <hr className="border-gray-500 mt-4 mb-4"></hr>
        </div>
      </div>
      <div className="flex flex-col ml-4">
        <SidebarButton>
          <i className="fa-solid fa-plus text-xs mr-2"></i>
          Create Team
        </SidebarButton>
        {teams.map((team, index) => (
          <SidebarButton key={index}>{team}</SidebarButton>
        ))}
      </div>
      <div className="absolute bottom-0 left-0 w-full p-2 pb-6">
        <hr className="border-gray-500 mt-4 mb-4"></hr>
        <SidebarButton>Create new workspace</SidebarButton>
        <SidebarButton>Change workspace</SidebarButton>
      </div>
    </div>
  );
}
