<script lang="ts">
  import { ListGames, LaunchGame } from "$lib/wailsjs/go/main/App.js";
  import { onMount } from "svelte";
  import toast from "svelte-french-toast";
  import type {FormEventHandler} from "svelte/elements";

  let games: any = [];

  $: filtered = games;

  function searchBarFilter(event: FormEventHandler<HTMLInputElement>) {
      const value = event.target.value.toLowerCase();
      filtered = games.filter((game: any) => game.Name.toLowerCase().includes(value));
  }

  onMount(async () => {
    games = await ListGames(); 

    toast.success("Games loaded", {
        position: "top-right",
        style: 'border-radius: 200px; background: #333; color: #fff;'
    });

    setInterval(async () => {
        games = await ListGames();
    }, 3000);
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
    <input type="text" placeholder="Pokémon White" class="input input-bordered w-full max-w-xs p-2 m-4" on:input={searchBarFilter} />

    <table class="table">
      <thead>
        <tr>
          <th>Icon</th>
          <th>Name</th> 
          <th>Developer</th>
          <th>Size</th> 
        </tr>
      </thead> 
      <tbody>
        {#each filtered as game}
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
