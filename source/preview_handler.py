import os
import json
import datetime
from source import configuration


class PreviewHandler:
    """
    Handles preview and dry-run functionality for the newsletter
    """
    
    def __init__(self):
        self.config = configuration.conf.preview
        self._ensure_output_directory()
    
    def _resolve_output_directory(self):
        """
        Resolve output directory based on environment and platform
        Returns the actual path where files should be saved
        """
        configured_path = self.config.output_directory
        
        # Check if we're in a Docker container by looking for /.dockerenv
        # This is the standard way to detect Docker environment
        is_docker = os.path.exists('/.dockerenv')
        
        if is_docker:
            # In Docker, use configured path as-is
            return configured_path
        else:
            # Not in Docker - handle Docker-style paths
            if configured_path.startswith('/app/config/'):
                # Convert Docker path to local relative path
                relative_path = configured_path.replace('/app/config/', './config/')
                return relative_path
            elif os.path.isabs(configured_path):
                # Absolute path - use as-is (works on all platforms)
                return configured_path
            else:
                # Relative path - use as-is (works on all platforms)
                return configured_path
    
    def _ensure_output_directory(self):
        """Create output directory if it doesn't exist"""
        if self.config.enabled:
            actual_path = self._resolve_output_directory()
            try:
                os.makedirs(actual_path, exist_ok=True)
                # Log the actual path being used
                configuration.logging.info(f"Preview directory: {os.path.abspath(actual_path)}")
            except Exception as e:
                configuration.logging.error(f"Failed to create preview directory '{actual_path}': {e}")
                raise
    
    def _generate_filename(self, suffix=""):
        """Generate filename with date/timestamp placeholders"""
        filename = self.config.output_filename
        now = datetime.datetime.now()
        
        # Replace placeholders
        filename = filename.replace("{date}", now.strftime("%Y-%m-%d"))
        filename = filename.replace("{time}", now.strftime("%H%M%S"))
        filename = filename.replace("{timestamp}", now.strftime("%Y%m%d_%H%M%S"))
        
        if suffix:
            name, ext = os.path.splitext(filename)
            filename = f"{name}_{suffix}{ext}"
            
        # Use resolved directory path
        output_dir = self._resolve_output_directory()
        return os.path.join(output_dir, filename)
    
    def _add_metadata_to_html(self, html_content, metadata):
        """Add metadata to HTML as comments"""
        if not self.config.include_metadata:
            return html_content
            
        metadata_comment = f"""<!--
=== JELLYFIN NEWSLETTER METADATA ===
Generated: {metadata['generation_timestamp']}
Mode: {metadata['mode']}
Movies Found: {metadata['stats']['movies_count']}
TV Episodes Found: {metadata['stats']['tv_episodes_count']}
Template Language: {metadata['template_language']}
SMTP Tested: {metadata.get('smtp_tested', 'N/A')}
=== END METADATA ===
-->
"""
        
        # Insert after DOCTYPE or at the beginning
        if '<!DOCTYPE' in html_content:
            parts = html_content.split('>', 1)
            if len(parts) == 2:
                return parts[0] + '>\n' + metadata_comment + '\n' + parts[1]
        
        return metadata_comment + '\n' + html_content
    
    def save_preview(self, html_content, metadata, mode="preview"):
        """
        Save HTML preview and optional JSON metadata
        
        Args:
            html_content (str): The generated HTML email content
            metadata (dict): Email generation metadata
            mode (str): "preview" or "dry-run"
            
        Returns:
            tuple: (html_file_path, json_file_path or None)
        """
        if not self.config.enabled:
            return None, None
            
        try:
            # Generate filenames
            html_file = self._generate_filename()
            
            # Add metadata to HTML
            if self.config.include_metadata:
                html_content = self._add_metadata_to_html(html_content, metadata)
            
            # Save HTML file
            with open(html_file, 'w', encoding='utf-8') as f:
                f.write(html_content)
            
            # Save JSON metadata if enabled
            json_file = None
            if self.config.save_email_data:
                json_file = self._generate_filename("data").replace('.html', '.json')
                with open(json_file, 'w', encoding='utf-8') as f:
                    json.dump(metadata, f, indent=2, ensure_ascii=False, default=str)
            
            # Log actual file locations
            configuration.logging.info(f"Preview saved: {os.path.abspath(html_file)}")
            if json_file:
                configuration.logging.info(f"Metadata saved: {os.path.abspath(json_file)}")
            
            return html_file, json_file
            
        except Exception as e:
            configuration.logging.error(f"Failed to save preview files: {e}")
            raise
    
    def get_metadata(self, movies, series, total_tv, total_movie, mode="preview", smtp_tested=False):
        """
        Generate metadata for the email
        
        Args:
            movies (dict): Movies data
            series (dict): Series data
            total_tv (int): Total TV episodes count
            total_movie (int): Total movie count
            mode (str): "preview" or "dry-run"
            smtp_tested (bool): Whether SMTP connection was tested
            
        Returns:
            dict: Email metadata
        """
        now = datetime.datetime.now()
        
        # Prepare movies data for JSON
        movies_list = []
        for movie_id, movie_data in movies.items():
            movies_list.append({
                "name": movie_data.get('name', 'Unknown'),
                "added_date": movie_data.get('created_on', '').split('T')[0],
                "tmdb_id": movie_data.get('tmdb_id', '')
            })
        
        # Prepare series data for JSON
        series_list = []
        for serie_id, serie_data in series.items():
            series_list.append({
                "series_name": serie_data.get('series_name', 'Unknown'),
                "seasons": serie_data.get('seasons', []),
                "episodes": serie_data.get('episodes', []),
                "added_date": serie_data.get('created_on', '').split('T')[0]
            })
        
        return {
            "generation_timestamp": now.isoformat(),
            "mode": mode,
            "smtp_tested": smtp_tested,
            "jellyfin_server": configuration.conf.email_template.jellyfin_url,
            "stats": {
                "movies_count": total_movie,
                "tv_episodes_count": total_tv,
                "total_email_size_kb": 0  # Will be calculated after HTML generation
            },
            "movies": movies_list,
            "tv_shows": series_list,
            "recipients": configuration.conf.recipients if mode == "dry-run" else ["preview-mode"],
            "template_language": configuration.conf.email_template.language,
            "configuration_hash": str(hash(str(configuration.conf.__dict__)))
        }
