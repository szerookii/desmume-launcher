<script lang="ts">
  import { ListGames, LaunchGame } from "$lib/wailsjs/go/main/App.js";
  import { onMount } from "svelte";
  import toast from "svelte-french-toast";

  let games: any = [];

  onMount(async () => {
    games = await ListGames();
  });

  async function onDoubleClick(game: any) {
      if(await LaunchGame(game)) {
        toast.success("Game launched", {
            position: "top-right",
            style: 'border-radius: 200px; background: #333; color: #fff;'
        });
      } else {
        toast.error("Failed to launch game", {
            position: "top-right",
            style: 'border-radius: 200px; background: #333; color: #fff;'
        });
      }
  }
</script>

<section class="min-h-screen scroll-smooth">
  <div class="overflow-x-auto">
    <table class="table">
      <thead>
        <tr>
          <th></th> 
          <th>Name</th> 
          <th>Developer</th>
          <th>Size</th> 
        </tr>
      </thead> 
      <tbody>
        {#each games as game}
        <tr on:dblclick={(() => onDoubleClick(game))}>
          <th><div class="avatar">
            <div class="w-16 h-16">
              <img src={"data:image/png;base64," + game.Base64EncodedIcon} alt="icon" />
            </div>
          </div></th>
          <td>{game.Name}</td>
          <td>{game.Developer}</td>
          <td>{game.Size}</td> 
        </tr>
        {/each}
      </tbody> 
    </table>
  </div>
</section>
