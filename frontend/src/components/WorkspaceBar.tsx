"use client";

import { useRouter } from "next/navigation";
import Link from "next/link";

export default function WorkspaceBar() {
  const currentWorkspace = "Workspace 1";
  const teams = ["Team 1", "Team 2", "Team 3"];
  const workspaces = ["Workspace 1", "Workspace 2", "Workspace 3"];
  const joinedWorkspaces = ["Joined Workspace 1", "Joined Workspace 2"];
  const router = useRouter();
  const handleLogout = () => {
    router.push("/login");
  };
  const createWorkspace = () => {
    router.push("/workspace-edit");
  };
  return (
    <div className="flex-column items-center justify-between p-4 border-r-1 border-gray-300 bg-gray-800 text-white w-60 h-screen">
      <div className="flex items-center justify-between p-2 pb-4 mb-4 border-b border-gray-600 text-xl font-bold">
        Username
        <i
          className="fa-solid fa-right-from-bracket text-xl"
          onClick={handleLogout}
        ></i>
      </div>
      <div className="flex items-center justify-between p-2 text-gray-200 text-sm font-normal">
        Your Workspaces
        <i
          className="fa-solid fa-square-plus text-xl"
          onClick={createWorkspace}
        ></i>
      </div>
      <ul className="flex-column items-start justify-start h-3/10 overflow-y-auto">
        {workspaces.map((workspace, index) => (
          <li key={index} className="p-2 cursor-pointer hover:underline">
            <Link href={"/workspace"}>{workspace}</Link>
            {workspace === currentWorkspace && (
              <ul className="ml-4 mt-2">
                {teams.map((team, tIdx) => (
                  <li
                    key={tIdx}
                    className="text-sm hover:bg-gray-700 cursor-pointer text-gray-400 pl-2 py-1"
                  >
                    <Link href={"/team"}>{team}</Link>
                  </li>
                ))}
              </ul>
            )}
          </li>
        ))}
      </ul>
      <div className="flex items-center justify-between p-2 text-gray-200 text-sm font-normal">
        Joined Workspaces
      </div>
      <ul className="flex-column items-start justify-start h-3/10 overflow-y-auto">
        {joinedWorkspaces.map((workspace, index) => (
          <li key={index} className="p-2 cursor-pointer hover:underline">
            <Link href={"/workspace"}>{workspace}</Link>
            {workspace === currentWorkspace && (
              <ul className="ml-4 mt-2">
                {teams.map((team, tIdx) => (
                  <li
                    key={tIdx}
                    className="text-xs hover:bg-gray-700 cursor-pointer text-gray-400 pl-2 py-1"
                  >
                    <Link href={"/team"}>{team}</Link>
                  </li>
                ))}
              </ul>
            )}
          </li>
        ))}
      </ul>
    </div>
  );
}
