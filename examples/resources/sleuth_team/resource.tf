resource "sleuth_team" "sampleteam" {
   name = "Sample team"
   members = [1] # use user IDs from GraphQL API
}