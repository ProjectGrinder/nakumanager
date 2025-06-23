"use client";

import { useRouter } from "next/navigation";

export default function WorkspaceInfo() {
  const workspace_owner = "John Doe";
  const admin_list = ["Admin 1", "Admin 2", "Admin 3"];
  const member_list = [
    "Member 1",
    "Member 2",
    "Member 3",
    "Member 1",
    "Member 2",
    "Member 3",
    "Member 1",
    "Member 2",
    "Member 3",
    "Member 1",
    "Member 2",
    "Member 3",
  ];
  const router = useRouter();
  const editWorkspace = () => {
    router.push("/workspace-edit");
  };
  return (
    <div className="flex flex-col p-10 text-white w-3/5">
      <div className="flex-row text-2xl font-bold mb-4">
        <span>Workspace Owner</span>
        <i
          className="fa-solid fa-gear text-2xl ml-10"
          onClick={editWorkspace}
        ></i>
      </div>
      <span className="text-lg font-normal mb-10">{workspace_owner}</span>
      <div className="flex flex-row justify-between mb-4">
        <div className="flex-column items-left justify-start w-100">
          <span className="font-bold text-xl">Admins</span>
          <ul className="flex-column items-start justify-start h-100 mt-6 text-base overflow-y-auto">
            {admin_list.map((admin, index) => (
              <li className="pt-2 pb-2" key={index}>
                {admin}
              </li>
            ))}
          </ul>
        </div>
        <div className="flex-column items-left justify-start w-100">
          <span className="font-bold text-xl">Members</span>
          <ul className="flex-column items-start justify-start h-100 mt-6 text-base overflow-y-auto">
            {member_list.map((member, index) => (
              <li className="pt-2 pb-2" key={index}>
                {member}
              </li>
            ))}
          </ul>
        </div>
      </div>
    </div>
  );
}
