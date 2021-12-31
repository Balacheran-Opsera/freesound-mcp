#
# Freesound is (c) MUSIC TECHNOLOGY GROUP, UNIVERSITAT POMPEU FABRA
#
# Freesound is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as
# published by the Free Software Foundation, either version 3 of the
# License, or (at your option) any later version.
#
# Freesound is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with this program.  If not, see <http://www.gnu.org/licenses/>.
#
# Authors:
#     See AUTHORS file.
#

from __future__ import absolute_import

from django import template

from sounds.models import Pack, Sound
from utils.frontend_handling import using_beastwhoosh

register = template.Library()


@register.inclusion_tag('sounds/display_pack.html', takes_context=True)
def display_pack(context, pack, size='small'):
    """This templatetag is used to display a pack with some randomly selected sound players.

    Args:
        context (django.template.Context): an object with contextual information for rendering a template. This
          argument is automatically added by Django when calling the templatetag inside a template.
        pack (int or Pack): pack ID or Pack object of the pack that will be shown. If no pack exists for the
          given ID, the display_pack.html will be rendered with empty HTML.
        size (str, optional): size of the "info" to display. This parameter only applies to BW interface.
          Must be one of ['small' (default), 'big']. Information about the contents of each
          size is given in the display_pack.html template code.

    Returns:
        dict: dictionary with the variables needed for rendering the pack with the display_pack.html template

    """
    if isinstance(pack, Pack):
        pack_obj = [pack]
    else:
        try:
            # use filter here instead of get because we don't want the query to be evaluated before rendering the
            # template as this would bypass the html cache in the template
            pack_obj = Pack.objects.select_related('user').filter(id=int(pack))
        except ValueError:
            # Invalid ID, we set pack_obj to empty list so "if pack" check in template returns False
            pack_obj = []
        except Pack.DoesNotExist:
            # Pack does not exist, we set pack_obj to empty list so "if pack" check in template returns False
            pack_obj = []

    request = context.get('request')
    if using_beastwhoosh(request):
        try:
            user_profile_locations = pack_obj[0].user.profile.locations()
        except IndexError:
            user_profile_locations = None
    else:
        user_profile_locations = None

    # Add 'request' to the returned context dictionary below so when the display_sound templatetag is called inside
    # display_pack templatetag it is given request in the context as well.
    return {
        'pack': pack_obj,
        'size': size,
        'user_profile_locations': user_profile_locations,
        'media_url': context['media_url'],
        'request': request
    }


@register.inclusion_tag('sounds/display_pack.html', takes_context=True)
def display_pack_big(context, pack):
    return display_pack(context, pack, size='big')
