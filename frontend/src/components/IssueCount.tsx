export default function IssueCount(issue_list: string[]) {
  return (
    <div className="h-8 w-8 rounded-2xl bg-gray-300 text-gray-700 text-base flex justify-center items-center">
      <span>{issue_list.length}</span>
    </div>
  );
}
